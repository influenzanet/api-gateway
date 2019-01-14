package structs

// Config is the structure corresponding to config.yaml that holds all static configuration data
type Config struct {
	URLAuthenticationService string `yaml:"url_authentication_service"`
	URLUserManagementService string `yaml:"url_user_management_service"`
}
