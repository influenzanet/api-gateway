package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const abortResponseKey = "abortResponse"

// AbortResponse marks a response that failed after its headers were sent, so
// the client sees an incomplete transfer instead of a truncated success.
func AbortResponse(c *gin.Context) {
	c.Set(abortResponseKey, true)
}

// AbortFailedResponses closes the connection of responses marked with
// AbortResponse. Register it before gin's Recovery, which would otherwise
// swallow http.ErrAbortHandler.
func AbortFailedResponses() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.GetBool(abortResponseKey) {
			panic(http.ErrAbortHandler)
		}
	}
}
