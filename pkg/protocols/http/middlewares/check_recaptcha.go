package middlewares

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type RecaptchaValidationResp struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

// ValidateToken reads the token from the request and validates it by contacting the authentication service
func CheckRecaptcha(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request

		tokens, ok := req.Header["Recaptcha-Token"]
		if !ok || len(tokens) < 1 || len(tokens[0]) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "recaptcha token missing"})
			c.Abort()
			return
		}

		token := tokens[0]

		resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
			"secret":   {secretKey},
			"response": {token},
		})
		if err != nil {
			log.Printf("unexpected error during recaptcha validation: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error during recaptcha validation"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("unexpected error during recaptcha validation: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error during recaptcha validation"})
			c.Abort()
			return
		}

		var parsedResp RecaptchaValidationResp
		err = json.Unmarshal(body, &parsedResp)
		if err != nil {
			log.Printf("unexpected error during recaptcha validation: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error during recaptcha validation"})
			c.Abort()
			return
		}

		if !parsedResp.Success {
			log.Printf("recaptcha validation failed: %v", parsedResp)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error during recaptcha validation"})
			c.Abort()
			return
		}
		c.Next()
	}
}
