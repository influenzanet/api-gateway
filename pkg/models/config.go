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
	SAMLConfig *SAMLConfig `yaml:"saml_config"`
}

type UseEndpoints struct {
	DeleteParticipantData bool
	SignupWithEmail       bool
	LoginWithExternalIDP  bool
}

type SAMLConfig struct {
	IDPUrl                   string `yaml:"idp_root_url"`
	SPRootUrl                string `yaml:"sp_root_url"`
	EntityID                 string `yaml:"entity_id"`
	MetaDataURL              string `yaml:"metadata_url"`
	SessionCertPath          string `yaml:"session_cert"`
	SessionKeyPath           string `yaml:"session_key"`
	TemplatePathLoginSuccess string `yaml:"templates_login_success"`
}
