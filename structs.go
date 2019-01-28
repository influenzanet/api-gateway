package main

// Config is the structure corresponding to config.yaml that holds all static configuration data
type Config struct {
	ServiceURL struct {
		Authentication string `yaml:"authentication"`
		UserManagement string `yaml:"user_management"`
	} `yaml:"service_url"`
}
