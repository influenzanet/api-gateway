package v1

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	mw "github.com/influenzanet/api-gateway/pkg/protocols/http/middlewares"
	"github.com/influenzanet/api-gateway/pkg/utils"
	"github.com/influenzanet/go-utils/pkg/api_types"
	studyAPI "github.com/influenzanet/study-service/pkg/api"
	"google.golang.org/grpc/status"

	"github.com/h2non/filetype"
)

const chunkSize = 64 * 1024 // 64 KiB

func (h *HttpEndpoints) AddRTRSpecificEndpoints(rg *gin.RouterGroup) {
	auth := rg.Group("/rtr")
	auth.GET("/code-validation", h.rtrCodeValidationEndpoint)
	auth.POST("/sync/:studyKey/:start/:end", mw.ExtractToken(), mw.ValidateToken(h.clients.UserManagement), h.uploadRTRSync)
}

func (h *HttpEndpoints) uploadRTRSync(c *gin.Context) {
	log.Println("File upload")
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	studyKey := c.Param("studyKey")
	start := c.Param("start")
	end := c.Param("end")

	file, err := c.FormFile("file")
	// file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(file.Filename)
	log.Println(file.Size)

	// Open file
	f, err := file.Open()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	// Get bytes
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content := buf.Bytes()

	stream, err := h.clients.StudyService.UploadParticipantFile(context.Background())
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	kind, _ := filetype.Match(content)
	log.Printf("%v", kind)
	log.Printf("%v", kind.MIME.Value)

	// Send infos
	req := &studyAPI.UploadParticipantFileReq{
		Data: &studyAPI.UploadParticipantFileReq_Info_{
			Info: &studyAPI.UploadParticipantFileReq_Info{
				Token:                token,
				StudyKey:             studyKey,
				VisibleToParticipant: true,
				FileType: &studyAPI.FileType{
					Type:    kind.MIME.Type,
					Subtype: kind.MIME.Subtype,
					Value:   kind.MIME.Value,
				},
				Participant: &studyAPI.UploadParticipantFileReq_Info_ProfileId{
					ProfileId: token.ProfilId,
				},
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return

	}

	for currentByte := 0; currentByte < len(content); currentByte += chunkSize {
		var currentChnk []byte
		if currentByte+chunkSize > len(content) {
			currentChnk = content[currentByte:]
		} else {
			currentChnk = content[currentByte : currentByte+chunkSize]
		}
		log.Println(len(currentChnk))
		req = &studyAPI.UploadParticipantFileReq{
			Data: &studyAPI.UploadParticipantFileReq_Chunk{
				Chunk: currentChnk,
			},
		}

		if err := stream.Send(req); err != nil {
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return

	}

	sReq := studyAPI.SubmitResponseReq{
		Token:     token,
		StudyKey:  c.Param("studyKey"),
		ProfileId: token.ProfilId,
		Response: &studyAPI.SurveyResponse{
			SubmittedAt: time.Now().Unix(),
			Key:         "rtrsyncfile",
			Responses: []*studyAPI.SurveyItemResponse{
				{Key: "file", Response: &studyAPI.ResponseItem{
					Key: "rg", Items: []*studyAPI.ResponseItem{
						{Key: "input", Value: reply.Id},
					},
				}},
				{Key: "start", Response: &studyAPI.ResponseItem{
					Key: "rg", Items: []*studyAPI.ResponseItem{
						{Key: "input", Value: start},
					},
				}},
				{Key: "end", Response: &studyAPI.ResponseItem{
					Key: "rg", Items: []*studyAPI.ResponseItem{
						{Key: "input", Value: end},
					},
				}},
			},
		},
	}
	_, err = h.clients.StudyService.SubmitResponse(context.Background(), &sReq)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	h.SendProtoAsJSON(c, http.StatusOK, reply)
}

func (h *HttpEndpoints) rtrCodeValidationEndpoint(c *gin.Context) {
	code, codeExists := c.GetQuery("code")

	if !codeExists || len(code) < 4 {
		time.Sleep(time.Second * 2)
		c.JSON(http.StatusBadRequest, gin.H{"error": "code error"})
		return
	}
	log.Printf("GROUP CODE: %s", code)

	for k, codes := range codeLists {
		for _, ec := range codes {
			if ec == code {
				c.JSON(http.StatusOK, gin.H{
					"group": k,
					"code":  code,
				})
				return
			}
		}
	}

	time.Sleep(time.Second * 2)
	c.JSON(http.StatusBadRequest, gin.H{"error": "code error"})
}

var codeLists map[string][]string = map[string][]string{
	"EG": {
		"eg08",
		"testeg",
		"i6ECbF",
		"FgPnim",
		"FT6HtP",
		"tKPypo",
		"UZWmzN",
		"H2SCQF",
		"QR856t",
		"Bjd4ew",
		"axnxYG",
		"WaAMxn",
		"yWQsFN",
		"83nrhm",
		"XAegyd",
		"mtFL6m",
		"T9LTqK",
		"bdByYS",
		"4AXm6r",
		"S5cjWK",
		"ooig4R",
		"rmR87N",
		"4Mma7f",
		"J4Lvqq",
		"ntHRTg",
		"LaSy76",
		"EPgLFD",
		"dhjxSN",
		"6MrRAP",
		"xc9tBq",
		"oU9sfV",
		"3UHjmU",
		"JshoCA",
		"BKggC5",
		"erZ4uw",
		"5vbV5D",
		"twDvBF",
		"WT8BZU",
		"dtfZyA",
		"5wMDfZ",
		"APS4Kc",
		"UzyDKf",
		"ZQvxGx",
		"YYQeNE",
		"zWfEGW",
		"CVv8QP",
		"kyV5Qn",
		"6mqduu",
		"4CUbUE",
		"HtUeWW",
		"QDd5TW",
		"xLn25j",
		"X89tPT",
		"5T9LZX",
		"GSTdf4",
		"Q8GGtj",
		"UZnvsv",
		"2qeYLw",
		"kNCHTk",
		"Gnit9Z",
		"CZTAWA",
		"LswGiV",
		"bcnm8E",
		"tvQDMV",
		"uP22Mn",
		"yGpaAs",
		"fzG4S5",
		"rWoe4T",
		"xyFNCy",
		"whYDkM",
		"puwmc9",
		"4CeDj6",
		"FJtvsK",
		"qCryq6",
		"RBDCwL",
		"UZYWqj",
		"QMWakZ",
		"HBQD7B",
		"yoUHEC",
		"2uJHzs",
		"uwVaD8",
		"nP6vTb",
		"5R5bpv",
		"3Vxa6Y",
		"Mb9NSZ",
		"zZ4mhN",
		"JAtYW5",
		"ibdxTq",
		"XJNyTF",
		"c43W7h",
		"Lqs8pk",
		"XbKLgS",
		"DG8ftE",
		"bnpHGe",
		"RuTEQ4",
		"Zc867f",
		"gQcpJu",
		"Put79F",
		"fvQjzB",
		"JmCGTr",
		"oqzpht",
		"miGUVf",
	},
	"EG09": {
		"eg09",
	},
	"PUB08": {
		"pub08",
	},
	"PUB09": {
		"pub09",
	},
	"RTR05": {
		"rtr05",
	},
	"RTR60": {
		"rtr60",
	},
	"TEST01": {
		"empty",
	},
}
