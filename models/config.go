package models

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
