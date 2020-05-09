package models

import (
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
)

// APIClients holds the service clients to the internal services
type APIClients struct {
	UserManagement umAPI.UserManagementApiClient
	// AuthService    api.AuthServiceApiClient
	// StudyService   api.StudyServiceApiClient
}
