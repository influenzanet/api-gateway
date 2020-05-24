package utils

import (
	messageAPI "github.com/influenzanet/messaging-service/pkg/api/messaging_service"
	studyAPI "github.com/influenzanet/study-service/pkg/api"
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
)

func ConvertTokenInfosForStudyAPI(t *umAPI.TokenInfos) *studyAPI.TokenInfos {
	if t == nil {
		return nil
	}
	return &studyAPI.TokenInfos{
		Id:               t.Id,
		InstanceId:       t.InstanceId,
		AccountConfirmed: t.AccountConfirmed,
		IssuedAt:         t.IssuedAt,
		Payload:          t.Payload,
		ProfilId:         t.ProfilId,
	}
}

func ConvertTokenInfosForMessageAPI(t *umAPI.TokenInfos) *messageAPI.TokenInfos {
	if t == nil {
		return nil
	}
	return &messageAPI.TokenInfos{
		Id:               t.Id,
		InstanceId:       t.InstanceId,
		AccountConfirmed: t.AccountConfirmed,
		IssuedAt:         t.IssuedAt,
		Payload:          t.Payload,
		ProfilId:         t.ProfilId,
	}
}
