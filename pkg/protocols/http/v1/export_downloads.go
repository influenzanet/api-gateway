package v1

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/coneno/logger"
	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/protocols/http/middlewares"
	"github.com/influenzanet/api-gateway/pkg/utils"
	studyAPI "github.com/influenzanet/study-service/pkg/api"
	"google.golang.org/grpc/status"
)

type responseExportStream interface {
	Recv() (*studyAPI.Chunk, error)
}

// serveResponseExport buffers the export by default. With ?stream=true chunks are
// sent as they arrive, a later upstream failure then aborts the connection.
func serveResponseExport(c *gin.Context, contentType, filename string,
	openStream func(context.Context) (responseExportStream, error),
) {
	stream, err := openStream(c.Request.Context())
	if err != nil {
		writeExportError(c, err)
		return
	}
	headers := map[string]string{"Content-Disposition": "attachment; filename=" + filename}
	if c.Query("stream") == "true" {
		streamResponseExport(c, stream, contentType, headers)
		return
	}

	content := []byte{}
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			writeExportError(c, err)
			return
		}
		content = append(content, chunk.Chunk...)
	}
	c.DataFromReader(http.StatusOK, int64(len(content)), contentType, bytes.NewReader(content), headers)
}

func streamResponseExport(c *gin.Context, stream responseExportStream, contentType string, headers map[string]string) {
	started := false
	start := func() {
		for key, value := range headers {
			c.Header(key, value)
		}
		c.Header("Content-Type", contentType)
		c.Status(http.StatusOK)
		started = true
	}
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			if !started {
				writeExportError(c, err)
			} else if c.Request.Context().Err() == nil {
				logger.Error.Printf("response export failed after download started: %v", err)
				middlewares.AbortResponse(c)
			}
			return
		}
		if !started {
			start()
		}
		if _, err := c.Writer.Write(chunk.Chunk); err != nil {
			return
		}
		c.Writer.Flush()
	}
	if !started {
		start()
		c.Writer.WriteHeaderNow()
	}
}

func writeExportError(c *gin.Context, err error) {
	st := status.Convert(err)
	c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
}
