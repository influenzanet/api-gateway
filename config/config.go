package config

import (
	"io/ioutil"
	"log"

	"github.com/Influenzanet/api-gateway/structs"
	yaml "gopkg.in/yaml.v2"
)

// Conf holds all static configuration information
var Conf structs.Config

// ReadConfig reads the config.yaml file and creates the config structure for all other packages to access
func ReadConfig() {
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal([]byte(data), &Conf)
	if err != nil {
		log.Fatal(err)
	}
}
