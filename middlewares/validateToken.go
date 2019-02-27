package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	auth_api "github.com/influenzanet/api/dist/go/auth-service"
	"google.golang.org/grpc/status"
)

// ValidateToken reads the token from the request and validates it by contacting the authentication service
func ValidateToken(authServiceClient auth_api.AuthServiceApiClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.MustGet("encodedToken").(string)
		req := auth_api.EncodedToken{
			Token: token,
		}
		parsedToken, err := authServiceClient.ValidateJWT(context.Background(), &req)
		if err != nil {
			st := status.Convert(err)
			log.Println(st.Message())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error during token validation"})
			c.Abort()
			return
		}
		c.Set("parsedToken", *parsedToken)

		c.Next()
	}
}
