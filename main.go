package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var client = &http.Client{}

func main() {
	ReadConfig()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	v1 := router.Group("/v1")

	InitUserEndpoints(v1)
	InitTokenEndpoints(v1)
	InitSurveyEndpoints(v1)

	log.Fatal(router.Run(":3000"))
}
