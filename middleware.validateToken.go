package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidateToken reads the token from the request and validates it by contacting the authentication service
func ValidateToken(conf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request

		var token string
		tokens, ok := req.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			token = tokens[0]
			token = strings.TrimPrefix(token, "Bearer ")
			if len(token) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "no Authorization token found"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no Authorization token found"})
			c.Abort()
			return
		}

		client := &http.Client{}

		req, err := http.NewRequest("GET", conf.URLAuthenticationService, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating reuest"})
			c.Abort()
			return
		}
		req.Header.Add("Authorization", strings.Join([]string{"Bearer", token}, " "))

		res, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating request"})
			c.Abort()
			return
		}
		defer res.Body.Close()

		rawBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating request"})
			c.Abort()
			return
		}
		var body map[string]string
		if err := json.Unmarshal(rawBody, &body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating request"})
			c.Abort()
			return
		}

		errValue, errExists := body["error"]
		if errExists {
			c.JSON(res.StatusCode, gin.H{"error": errValue})
			c.Abort()
			return
		}

		tokenValue, tokenExists := body["token"]
		if !tokenExists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error authenticating request"})
			c.Abort()
			return
		}

		c.Set("token", tokenValue)

		c.Next()
	}
}
