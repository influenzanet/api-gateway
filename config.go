package main

import "os"

func initConfig() {
	conf.DebugMode = os.Getenv("DEBUG_MODE") == "true"
	conf.Port = os.Getenv("GATEWAY_LISTEN_PORT")
	conf.ServiceURLs.Authentication = os.Getenv("ADDR_AUTH_SERVICE")
	conf.ServiceURLs.UserManagement = os.Getenv("ADDR_USER_MANAGEMENT_SERVICE")
}
