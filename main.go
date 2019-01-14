package main

import (
	"log"

	"github.com/Influenzanet/api-gateway/endpoints"
	"github.com/derekparker/delve/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ReadConfig()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	v1 := router.Group("/v1")

	endpoints.InitEndpoints(v1)

	log.Fatal(router.Run(":3000"))
}
