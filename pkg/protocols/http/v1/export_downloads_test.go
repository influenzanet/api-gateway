package v1

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/models"
	"github.com/influenzanet/api-gateway/pkg/protocols/http/middlewares"
	"github.com/influenzanet/go-utils/pkg/api_types"
	studyAPI "github.com/influenzanet/study-service/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeExportStream struct {
	grpc.ClientStream
	chunks [][]byte
	err    error
}

func (s *fakeExportStream) Recv() (*studyAPI.Chunk, error) {
	if len(s.chunks) > 0 {
		chunk := s.chunks[0]
		s.chunks = s.chunks[1:]
		return &studyAPI.Chunk{Chunk: chunk}, nil
	}
	if s.err != nil {
		return nil, s.err
	}
	return nil, io.EOF
}

type exportStudyClient struct {
	studyAPI.StudyServiceApiClient
	open func(*studyAPI.ResponseExportQuery) (*fakeExportStream, error)
}

func (s *exportStudyClient) GetResponsesWideFormatCSV(_ context.Context, req *studyAPI.ResponseExportQuery, _ ...grpc.CallOption) (studyAPI.StudyServiceApi_GetResponsesWideFormatCSVClient, error) {
	return s.open(req)
}
func (s *exportStudyClient) GetResponsesLongFormatCSV(_ context.Context, req *studyAPI.ResponseExportQuery, _ ...grpc.CallOption) (studyAPI.StudyServiceApi_GetResponsesLongFormatCSVClient, error) {
	return s.open(req)
}
func (s *exportStudyClient) GetResponsesFlatJSON(_ context.Context, req *studyAPI.ResponseExportQuery, _ ...grpc.CallOption) (studyAPI.StudyServiceApi_GetResponsesFlatJSONClient, error) {
	return s.open(req)
}

// exportTestServer mirrors the management API middleware order.
func exportTestServer(t *testing.T, open func(*studyAPI.ResponseExportQuery) (*fakeExportStream, error)) *httptest.Server {
	t.Helper()
	h := NewHTTPHandler(&models.APIClients{StudyService: &exportStudyClient{open: open}}, models.UseEndpoints{}, nil)
	router := gin.New()
	router.Use(middlewares.AbortFailedResponses(), gin.Recovery())
	router.Use(func(c *gin.Context) { c.Set("validatedToken", &api_types.TokenInfos{}) })
	path := "/data/:studyKey/survey/:surveyKey/response"
	router.GET(path, h.getResponseWideFormatCSV)
	router.GET(path+"/long-format", h.getResponseLongFormatCSV)
	router.GET(path+"/json", h.getResponseFlatJSON)
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	return server
}

var exportFormats = []struct{ suffix, contentType, filename string }{
	{"", "text/csv", "study_survey.csv"},
	{"/long-format", "text/csv", "study_survey.csv"},
	{"/json", "application/json", "study_survey.json"},
}

func exportScenario(scenario string) (*fakeExportStream, error) {
	switch scenario {
	case "open error":
		return nil, status.Error(codes.Unavailable, "upstream unavailable")
	case "error before data":
		return &fakeExportStream{err: status.Error(codes.PermissionDenied, "upstream denied")}, nil
	case "error after data":
		return &fakeExportStream{chunks: [][]byte{[]byte("abc")}, err: errors.New("upstream failed")}, nil
	case "empty":
		return &fakeExportStream{}, nil
	}
	return &fakeExportStream{chunks: [][]byte{[]byte("abc"), []byte("def")}}, nil
}

func TestResponseExports(t *testing.T) {
	for _, format := range exportFormats {
		for _, query := range []string{"", "&stream=true"} {
			for _, tc := range []struct {
				scenario string
				code     int
				body     string
			}{
				{"success", http.StatusOK, "abcdef"},
				{"empty", http.StatusOK, ""},
				{"open error", http.StatusServiceUnavailable, `{"error":"upstream unavailable"}`},
				{"error before data", http.StatusUnauthorized, `{"error":"upstream denied"}`},
			} {
				t.Run(format.suffix+query+"/"+tc.scenario, func(t *testing.T) {
					server := exportTestServer(t, func(req *studyAPI.ResponseExportQuery) (*fakeExportStream, error) {
						if req.StudyKey != "study" || req.SurveyKey != "survey" || req.Page != 2 || req.Separator != "." {
							t.Errorf("export parameters changed: %+v", req)
						}
						return exportScenario(tc.scenario)
					})
					res, err := http.Get(server.URL + "/data/study/survey/survey/response" + format.suffix + "?page=2&sep=." + query)
					if err != nil {
						t.Fatal(err)
					}
					defer res.Body.Close()
					body, err := io.ReadAll(res.Body)
					if err != nil || res.StatusCode != tc.code || string(body) != tc.body {
						t.Fatalf("unexpected response: code=%d body=%q err=%v", res.StatusCode, body, err)
					}
					disposition := res.Header.Get("Content-Disposition")
					if tc.code != http.StatusOK {
						if disposition != "" {
							t.Fatalf("error response has download headers: %v", res.Header)
						}
						return
					}
					if disposition != "attachment; filename="+format.filename || res.Header.Get("Content-Type") != format.contentType {
						t.Fatalf("download headers changed: %v", res.Header)
					}
					// Buffered downloads declare their length; streamed ones are chunked.
					wantLength := int64(len(tc.body))
					if query != "" && tc.body != "" {
						wantLength = -1
					}
					if res.ContentLength != wantLength {
						t.Fatalf("unexpected content length %d", res.ContentLength)
					}
				})
			}
		}
	}
}

func TestBufferedExportReportsLateErrors(t *testing.T) {
	server := exportTestServer(t, func(*studyAPI.ResponseExportQuery) (*fakeExportStream, error) {
		return exportScenario("error after data")
	})
	res, err := http.Get(server.URL + "/data/study/survey/survey/response")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusInternalServerError || string(body) != `{"error":"upstream failed"}` {
		t.Fatalf("unexpected response: code=%d body=%q", res.StatusCode, body)
	}
}

func TestStreamedExportAbortsOnLateErrors(t *testing.T) {
	server := exportTestServer(t, func(*studyAPI.ResponseExportQuery) (*fakeExportStream, error) {
		return exportScenario("error after data")
	})
	res, err := http.Get(server.URL + "/data/study/survey/survey/response?stream=true")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK || string(body) != "abc" || !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated stream was not aborted: code=%d body=%q err=%v", res.StatusCode, body, err)
	}
}
