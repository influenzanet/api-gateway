package v1

import (
	"github.com/gin-gonic/gin"
	mw "github.com/influenzanet/api-gateway/pkg/protocols/http/middlewares"
)

func (h *HttpEndpoints) AddUserManagementParticipantAPI(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/resend-verification-code", mw.RequirePayload(), h.resendVerificationCodeHandl)
	auth.POST("/get-verification-code-with-token", mw.RequirePayload(), h.getVerificationCodeWithTokenHandl)

	auth.POST("/login-with-email", mw.RequirePayload(), h.loginWithEmailAsParticipantHandl)
	auth.POST("/signup-with-email", mw.RequirePayload(), h.signupWithEmailHandl)
	auth.POST("/switch-profile", mw.ExtractToken(), mw.ValidateToken(h.clients.UserManagement), mw.RequirePayload(), h.switchProfileHandl)
	auth.POST("/renew-token", mw.ExtractToken(), mw.RequirePayload(), h.tokenRenewHandl)

	user := rg.Group("/user")
	user.Use(mw.ExtractToken())
	user.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		user.GET("", h.getUserHandl)
		// userToken.GET("/:id", getUserHandl)
		user.POST("/change-password", mw.RequirePayload(), h.userPasswordChangeHandl)
		user.POST("/change-account-email", mw.RequirePayload(), h.changeAccountEmailHandl)
		user.POST("/revoke-refresh-tokens", h.revokeRefreshTokensHandl)
		user.POST("/set-language", mw.RequirePayload(), h.userSetPreferredLanguageHandl)
		user.POST("/delete", mw.RequirePayload(), h.deleteAccountHandl)

		user.POST("/profile/save", mw.CheckAccountConfirmed(), mw.RequirePayload(), h.saveProfileHandl)
		user.POST("/profile/remove", mw.RequirePayload(), h.removeProfileHandl)

		user.POST("/resend-verification-message", mw.RequirePayload(), h.resendContanctVerificationEmailHandl)
		user.POST("/contact-preferences", mw.RequirePayload(), h.userUpdateContactPreferencesHandl)
		user.POST("/contact/add-email", mw.CheckAccountConfirmed(), mw.RequirePayload(), h.userAddEmailHandl)
		user.POST("/contact/remove-email", mw.RequirePayload(), h.userRemoveEmailHandl)
	}

	unAuthUser := rg.Group("/user")
	{
		unAuthUser.POST("/password-reset/initiate", mw.RequirePayload(), h.initiatePasswordResetHandl)
		unAuthUser.POST("/password-reset/get-infos", mw.RequirePayload(), h.getInfosForPasswordResetHandl)
		unAuthUser.POST("/password-reset/reset-with", mw.RequirePayload(), h.passwordResetHandl)

		unAuthUser.POST("/contact-verification", mw.RequirePayload(), h.verifyUserContactHandl)
		unAuthUser.GET("/unsubscribe-newsletter", h.unsubscribeNewsletterHandl)
	}
}

func (h *HttpEndpoints) AddUserManagementAdminAPI(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/login-with-email", mw.RequirePayload(), h.loginWithEmailForManagementHandl)

	user := rg.Group("/user")
	user.Use(mw.ExtractToken())
	user.Use(mw.ValidateToken(h.clients.UserManagement))
	user.Use(mw.CheckAccountConfirmed())
	{
		user.POST("/", mw.RequirePayload(), h.createUserHandl)
		user.POST("/migrate", mw.RequirePayload(), h.migrateUserHandl)
		user.GET("/", h.findNonParticipantUsersHandl)
		user.POST("/add-role", mw.RequirePayload(), h.userAddRoleHandl)
		user.POST("/remove-role", mw.RequirePayload(), h.userRemoveRoleHandl)
	}
}
