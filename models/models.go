package models

import api "github.com/influenzanet/api-gateway/api"

// APIClients holds the service clients to the internal services
type APIClients struct {
	UserManagement api.UserManagementApiClient
	AuthService    api.AuthServiceApiClient
	StudyService   api.StudyServiceApiClient
}
