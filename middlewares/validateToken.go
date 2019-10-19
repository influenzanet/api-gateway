package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/influenzanet/api-gateway/api"
	"google.golang.org/grpc/status"
)

// ValidateToken reads the token from the request and validates it by contacting the authentication service
func ValidateToken(authClient api.AuthServiceApiClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.MustGet("encodedToken").(string)
		parsedToken, err := authClient.ValidateJWT(context.Background(), &api.JWTRequest{
			Token: token,
		})
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
