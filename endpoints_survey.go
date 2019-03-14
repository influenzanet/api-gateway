package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitSurveyEndpoints creates all API routes on the supplied RouterGroup
func InitSurveyEndpoints(rg *gin.RouterGroup) {
	survey := rg.Group("/survey") // TODO example endpoints
	{
		survey.POST("/submit", surveySubmitHandl)
		survey.POST("/update", surveyUpdateHandl)
		survey.POST("/get", surveyGetHandl)
		survey.POST("/get-all", surveyGetAllHandl)
	}
}

func surveySubmitHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func surveyUpdateHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func surveyGetHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func surveyGetAllHandl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
