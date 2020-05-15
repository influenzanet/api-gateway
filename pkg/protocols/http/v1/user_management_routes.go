package v1

import (
	"github.com/gin-gonic/gin"
	mw "github.com/influenzanet/api-gateway/pkg/protocols/http/middlewares"
)

func (h *HttpEndpoints) AddUserManagementParticipantAPI(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/login-with-email", mw.RequirePayload(), h.loginWithEmailAsParticipantHandl)
	auth.POST("/signup-with-email", mw.RequirePayload(), h.signupWithEmailHandl)
	auth.POST("/switch-profile", mw.ExtractToken(), mw.ValidateToken(h.clients.UserManagement), mw.RequirePayload(), h.switchProfileHandl)
	auth.POST("/renew-token", mw.ExtractToken(), mw.RequirePayload(), h.tokenRenewHandl)

	user := rg.Group("/user")
	user.Use(mw.ExtractToken())
	user.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		user.GET("/", h.getUserHandl)
		// userToken.GET("/:id", getUserHandl)
		user.POST("/change-password", mw.RequirePayload(), h.userPasswordChangeHandl)
		user.POST("/change-account-email", mw.RequirePayload(), h.changeAccountEmailHandl)
		user.POST("/revoke-refresh-tokens", h.revokeRefreshTokensHandl)
		user.POST("/set-language", mw.RequirePayload(), h.userSetPreferredLanguageHandl)
		user.POST("/delete", mw.RequirePayload(), h.deleteAccountHandl)

		user.POST("/profile/save", mw.RequirePayload(), h.saveProfileHandl)
		user.POST("/profile/remove", mw.RequirePayload(), h.removeProfileHandl)

		user.POST("/contact-preferences", mw.RequirePayload(), h.userUpdateContactPreferencesHandl)
		user.POST("/contact/add-email", mw.RequirePayload(), h.userAddEmailHandl)
		user.POST("/contact/remove-email", mw.RequirePayload(), h.userRemoveEmailHandl)
	}
}

func (h *HttpEndpoints) AddUserManagementAdminAPI(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/login-with-email", mw.RequirePayload(), h.loginWithEmailForManagementHandl)

	user := rg.Group("/user")
	user.Use(mw.ExtractToken())
	user.Use(mw.ValidateToken(h.clients.UserManagement))
	{
		user.POST("/", mw.RequirePayload(), h.createUserHandl)
		user.GET("/", h.findNonParticipantUsersHandl)
		user.POST("/add-role", mw.RequirePayload(), h.userAddRoleHandl)
		user.POST("/remove-role", mw.RequirePayload(), h.userRemoveRoleHandl)
	}
}
