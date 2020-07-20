package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/go-utils/pkg/api_types"
)

func CheckAccountConfirmed() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.MustGet("validatedToken").(*api_types.TokenInfos)

		if !token.AccountConfirmed {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "account not confirmed yet"})
			return
		}
		c.Next()
	}
}
