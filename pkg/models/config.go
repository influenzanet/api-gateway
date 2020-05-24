package models

// Config is the structure corresponding to config.yaml that holds all static configuration data
type Config struct {
	DebugMode   bool
	Port        string
	ServiceURLs struct {
		UserManagement   string `yaml:"user_management"`
		StudyService     string `yaml:"study_service"`
		MessagingService string `yaml:"messaging_service"`
	} `yaml:"service_urls"`
}
