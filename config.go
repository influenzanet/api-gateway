package main

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Conf holds all static configuration information
var conf Config

// ReadConfig reads the config.yaml file and creates the config structure for all other packages to access
func ReadConfig() {
	file := os.Getenv("CONFIG_FILE")
	if file == "" {
		file = "./config.yaml"
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		log.Fatal(err)
	}
}
