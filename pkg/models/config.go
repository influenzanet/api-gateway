package models

// Config is the structure corresponding to config.yaml that holds all static configuration data
type Config struct {
	DebugMode    bool
	AllowOrigins []string
	Port         string
	UseEndpoints UseEndpoints
	ServiceURLs  struct {
		UserManagement   string `yaml:"user_management"`
		StudyService     string `yaml:"study_service"`
		MessagingService string `yaml:"messaging_service"`
	} `yaml:"service_urls"`
}

type UseEndpoints struct {
	DeleteParticipantData bool
	SignupWithEmail       bool
}
