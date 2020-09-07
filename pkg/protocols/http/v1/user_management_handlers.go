package v1

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/influenzanet/api-gateway/pkg/utils"
	"github.com/influenzanet/go-utils/pkg/api_types"
	"github.com/influenzanet/go-utils/pkg/constants"
	"google.golang.org/grpc/status"

	studyAPI "github.com/influenzanet/study-service/pkg/api"
	umAPI "github.com/influenzanet/user-management-service/pkg/api"
)

func (h *HttpEndpoints) loginWithEmailAsParticipantHandl(c *gin.Context) {
	var req umAPI.LoginWithEmailMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.AsParticipant = true

	token, err := h.clients.UserManagement.LoginWithEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, token)
}

func (h *HttpEndpoints) resendVerificationCodeHandl(c *gin.Context) {
	var req umAPI.SendVerificationCodeReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.clients.UserManagement.SendVerificationCode(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getVerificationCodeWithTokenHandl(c *gin.Context) {
	var req umAPI.AutoValidateReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.clients.UserManagement.AutoValidateTempToken(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) loginWithEmailForManagementHandl(c *gin.Context) {
	var req umAPI.LoginWithEmailMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.AsParticipant = false

	resp, err := h.clients.UserManagement.LoginWithEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) signupWithEmailHandl(c *gin.Context) {
	var req umAPI.SignupWithEmailMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.clients.UserManagement.SignupWithEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) switchProfileHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	var req umAPI.SwitchProfileRequest
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.SwitchProfile(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) tokenRenewHandl(c *gin.Context) {
	var req umAPI.RefreshJWTRequest
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.AccessToken = c.MustGet("encodedToken").(string)
	token, err := h.clients.UserManagement.RenewJWT(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, token)
}

func (h *HttpEndpoints) userPasswordChangeHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.PasswordChangeMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token

	resp, err := h.clients.UserManagement.ChangePassword(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) changeAccountEmailHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.EmailChangeMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.ChangeAccountIDEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) userSetPreferredLanguageHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.LanguageChangeMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.ChangePreferredLanguage(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getUserHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	userID := c.Param("id")

	userRefReq := &umAPI.UserReference{
		Token:  token,
		UserId: userID,
	}

	resp, err := h.clients.UserManagement.GetUser(context.Background(), userRefReq)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) deleteAccountHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.UserReference
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.DeleteAccount(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		log.Println(st.Message())
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) saveProfileHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.ProfileRequest
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.SaveProfile(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) removeProfileHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.ProfileRequest
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.RemoveProfile(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) userUpdateContactPreferencesHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.ContactPreferencesMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.UpdateContactPreferences(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) userAddEmailHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.ContactInfoMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.AddEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) userRemoveEmailHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.ContactInfoMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.RemoveEmail(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) createUserHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.CreateUserReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.CreateUser(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

type MigrateUserReq struct {
	AccountID         string   `json:"accountId"`
	OldParticipantID  string   `json:"oldParticipantID"`
	InitialPassword   string   `json:"initialPassword"`
	PreferredLanguage string   `json:"preferredLanguage"`
	Studies           []string `json:"studies"`
	Use2FA            bool     `json:"use2FA"`
}

func (h *HttpEndpoints) migrateUserHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req MigrateUserReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	cuReq := umAPI.CreateUserReq{
		AccountId:         req.AccountID,
		InitialPassword:   req.InitialPassword,
		PreferredLanguage: req.InitialPassword,
		Roles:             []string{constants.USER_ROLE_PARTICIPANT},
		Token:             token,
		Use_2Fa:           req.Use2FA,
	}
	newUser, err := h.clients.UserManagement.CreateUser(context.Background(), &cuReq)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}

	newUserToken := &api_types.TokenInfos{
		Id:               newUser.Id,
		AccountConfirmed: true,
		InstanceId:       token.InstanceId,
		ProfilId:         newUser.Profiles[0].Id,
	}
	for _, studyKey := range req.Studies {
		// enter studies:
		_, err = h.clients.StudyService.EnterStudy(context.TODO(), &studyAPI.EnterStudyRequest{
			Token:    newUserToken,
			StudyKey: studyKey,
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}

		// submit migration survey
		_, err = h.clients.StudyService.SubmitResponse(context.TODO(), &studyAPI.SubmitResponseReq{
			Token:    newUserToken,
			StudyKey: studyKey,
			Response: &studyAPI.SurveyResponse{
				Key:         "migration",
				SubmittedAt: time.Now().Unix(),
				Responses: []*studyAPI.SurveyItemResponse{
					{Key: "migration.OldID", Response: &studyAPI.ResponseItem{Key: "rg", Items: []*studyAPI.ResponseItem{
						{Key: "ic", Value: req.OldParticipantID},
					}}},
				},
				Context: map[string]string{
					"engineVersion": "-",
				},
			},
		})
		if err != nil {
			st := status.Convert(err)
			c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
			return
		}
	}

	h.SendProtoAsJSON(c, http.StatusOK, newUser)
}

func (h *HttpEndpoints) findNonParticipantUsersHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.FindNonParticipantUsersMsg
	req.Token = token
	resp, err := h.clients.UserManagement.FindNonParticipantUsers(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) userAddRoleHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.RoleMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.AddRoleForUser(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) userRemoveRoleHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.RoleMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.RemoveRoleForUser(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) revokeRefreshTokensHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)

	var req umAPI.RevokeRefreshTokensReq
	req.Token = token
	resp, err := h.clients.UserManagement.RevokeAllRefreshTokens(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) initiatePasswordResetHandl(c *gin.Context) {
	var req umAPI.InitiateResetPasswordMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.clients.UserManagement.InitiatePasswordReset(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) getInfosForPasswordResetHandl(c *gin.Context) {
	var req umAPI.GetInfosForResetPasswordMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.clients.UserManagement.GetInfosForPasswordReset(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) passwordResetHandl(c *gin.Context) {
	var req umAPI.ResetPasswordMsg
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.clients.UserManagement.ResetPassword(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) verifyUserContactHandl(c *gin.Context) {
	var req umAPI.TempToken
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.clients.UserManagement.VerifyContact(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) resendContanctVerificationEmailHandl(c *gin.Context) {
	token := c.MustGet("validatedToken").(*api_types.TokenInfos)
	var req umAPI.ResendContactVerificationReq
	if err := h.JsonToProto(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Token = token
	resp, err := h.clients.UserManagement.ResendContactVerification(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}

func (h *HttpEndpoints) unsubscribeNewsletterHandl(c *gin.Context) {
	var req umAPI.TempToken
	req.Token = c.DefaultQuery("token", "")
	resp, err := h.clients.UserManagement.UseUnsubscribeToken(context.Background(), &req)
	if err != nil {
		st := status.Convert(err)
		c.JSON(utils.GRPCStatusToHTTP(st.Code()), gin.H{"error": st.Message()})
		return
	}
	h.SendProtoAsJSON(c, http.StatusOK, resp)
}
