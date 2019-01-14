package endpoints

import (
	"io/ioutil"
	"net/http"

	"github.com/Influenzanet/api-gateway/config"
	"github.com/gin-gonic/gin"
)

var client = &http.Client{}

func userLoginHandl(c *gin.Context) {
	req, err := http.NewRequest("POST", config.Conf.URLAuthenticationService+config.Conf.AuthenticationLogin, c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(res.StatusCode, rawBody)
}

func userSignupHandl(c *gin.Context) {
	req, err := http.NewRequest("POST", config.Conf.URLUserManagementService+config.Conf.AuthenticationSignup, c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(res.StatusCode, rawBody)
}

func userPasswordChangeHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func userEmailVerifyHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
