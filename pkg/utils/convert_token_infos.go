package utils

import (
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
