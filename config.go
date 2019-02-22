package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Conf holds all static configuration information
var conf Config

// ReadConfig reads the config.yaml file and creates the config structure for all other packages to access
func ReadConfig() {
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		log.Fatal(err)
	}
}
