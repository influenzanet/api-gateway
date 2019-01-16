package main

// Config is the structure corresponding to config.yaml that holds all static configuration data
type Config struct {
	URLAuthenticationService string `yaml:"authentication"`
	AuthenticationLogin      string `yaml:"authentication_login"`
	AuthenticationSignup     string `yaml:"authentication_signup"`
	URLUserManagementService string `yaml:"user_management"`
	UserManagementLogin      string `yaml:"user_management_login"`
	UserManagementSignup     string `yaml:"user_management_signup"`
}
