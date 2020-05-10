package models

import (
	studyAPI "github.com/influenzanet/study-service/pkg/api"
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
)

// APIClients holds the service clients to the internal services
type APIClients struct {
	UserManagement umAPI.UserManagementApiClient
	StudyService   studyAPI.StudyServiceApiClient
}
