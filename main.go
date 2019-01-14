package main

import (
	"io/ioutil"
	"log"

	"github.com/Influenzanet/api-gateway/endpoints"
	"github.com/Influenzanet/api-gateway/structs"

	"github.com/gin-gonic/gin"
	yaml "gopkg.in/yaml.v2"
)

var conf structs.Config

func readConfig() {
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	readConfig()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	v1 := router.Group("/v1")

	endpoints.InitEndpoints(v1, conf)

	log.Fatal(router.Run(":3000"))
}
