package models

import (
	messageAPI "github.com/influenzanet/messaging-service/pkg/api/messaging_service"
	studyAPI "github.com/influenzanet/study-service/pkg/api"
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
)

// APIClients holds the service clients to the internal services
type APIClients struct {
	UserManagement   umAPI.UserManagementApiClient
	StudyService     studyAPI.StudyServiceApiClient
	MessagingService messageAPI.MessagingServiceApiClient
}
