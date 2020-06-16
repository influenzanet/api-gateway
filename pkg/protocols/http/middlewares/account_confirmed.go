package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
)

func CheckAccountConfirmed() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.MustGet("validatedToken").(*umAPI.TokenInfos)

		if !token.AccountConfirmed {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "account not confirmed yet"})
			return
		}
		c.Next()
	}
}
