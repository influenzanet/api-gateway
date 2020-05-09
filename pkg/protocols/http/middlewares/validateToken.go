package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
	"google.golang.org/grpc/status"
)

// ValidateToken reads the token from the request and validates it by contacting the authentication service
func ValidateToken(authClient umAPI.UserManagementApiClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.MustGet("encodedToken").(string)
		parsedToken, err := authClient.ValidateJWT(context.Background(), &umAPI.JWTRequest{
			Token: token,
		})
		if err != nil {
			st := status.Convert(err)
			log.Println(st.Message())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error during token validation"})
			c.Abort()
			return
		}
		c.Set("validatedToken", parsedToken)
		c.Next()
	}
}
