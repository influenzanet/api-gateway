package middlewares

import (
	"context"
	"log"
	"net/http"

	auth_api "github.com/Influenzanet/api/dist/go/auth-service"
	"github.com/gin-gonic/gin"
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
