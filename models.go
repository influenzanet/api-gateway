package main

import api "github.com/influenzanet/api-gateway/api"

// Config is the structure corresponding to config.yaml that holds all static configuration data
type Config struct {
	DebugMode   bool
	Port        string
	ServiceURLs struct {
		Authentication string `yaml:"authentication"`
		UserManagement string `yaml:"user_management"`
		StudyService   string `yaml:"study_service"`
	} `yaml:"service_urls"`
}

// APIClients holds the service clients to the internal services
type APIClients struct {
	userManagement api.UserManagementApiClient
	authService    api.AuthServiceApiClient
	studyService   api.StudyServiceApiClient
}
