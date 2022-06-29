package middlewares

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/coneno/logger"
	"github.com/gin-gonic/gin"
)

const (
	ENV_USE_RECAPTCHA_FALLBACK      = "USE_RECAPTCHA"
	ENV_USE_RECAPTCHA_FOR_PREFIX    = "USE_RECAPTCHA_FOR_"
	ENV_RECAPTCHA_SECRET_FALLBACK   = "RECAPTCHA_SECRET"
	ENV_RECAPTCHA_SECRET_FOR_PREFIX = "RECAPTCHA_SECRET_FOR_"
)

type RecaptchaValidationResp struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func getInstanceID(header http.Header) (string, bool) {
	instanceID, ok := header["Instance-Id"]
	if !ok || len(instanceID) < 1 || len(instanceID[0]) < 1 {
		return "", false
	}
	return instanceID[0], true
}

func getRecaptchaConfig(hasInstanceID bool, instanceID string) (bool, string) {
	if !hasInstanceID {
		return os.Getenv(ENV_USE_RECAPTCHA_FALLBACK) == "true", os.Getenv(ENV_RECAPTCHA_SECRET_FALLBACK)
	}

	useRecaptcha := false
	secretKey := ""

	useRecaptchaValueForInstance, ok1 := os.LookupEnv(ENV_USE_RECAPTCHA_FOR_PREFIX + strings.ToUpper(instanceID))
	if ok1 {
		useRecaptcha = useRecaptchaValueForInstance == "true"
	} else {
		useRecaptcha = os.Getenv(ENV_USE_RECAPTCHA_FALLBACK) == "true"
	}

	secretKeyValueForInstance, ok2 := os.LookupEnv(ENV_RECAPTCHA_SECRET_FOR_PREFIX + strings.ToUpper(instanceID))
	if ok2 {
		secretKey = secretKeyValueForInstance
	} else {
		secretKey = os.Getenv(ENV_RECAPTCHA_SECRET_FALLBACK)
	}
	return useRecaptcha, secretKey
}

// ValidateToken reads the token from the request and validates it by contacting the authentication service
func CheckRecaptcha() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request

		instanceID, hasInstanceID := getInstanceID(req.Header)
		useRecaptcha, secretKey := getRecaptchaConfig(hasInstanceID, instanceID)
		if !useRecaptcha || len(secretKey) < 1 {
			c.Next()
			return
		}

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
			logger.Error.Printf("unexpected error during recaptcha validation: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error during recaptcha validation"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error.Printf("unexpected error during recaptcha validation: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error during recaptcha validation"})
			c.Abort()
			return
		}

		var parsedResp RecaptchaValidationResp
		err = json.Unmarshal(body, &parsedResp)
		if err != nil {
			logger.Error.Printf("unexpected error during recaptcha validation: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error during recaptcha validation"})
			c.Abort()
			return
		}

		if !parsedResp.Success {
			logger.Error.Printf("recaptcha validation failed: %v", parsedResp)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error during recaptcha validation"})
			c.Abort()
			return
		}
		c.Next()
	}
}
