package endpoints

import (
	"github.com/Influenzanet/api-gateway/config"
	"github.com/Influenzanet/api-gateway/middleware"
	"github.com/gin-gonic/gin"
)

// InitEndpoints creates all API routes on the supplied RouterGroup
func InitEndpoints(rg *gin.RouterGroup) {
	user := rg.Group("/user")
	user.Use(middleware.RequirePayload())
	{
		user.POST("/login", userLoginHandl)
		user.POST("/signup", userSignupHandl)
	}
	userToken := rg.Group("/user")
	userToken.Use(middleware.ValidateToken(config.Conf))
	userToken.Use(middleware.RequirePayload())
	{
		user.POST("/change-password", userPasswordChangeHandl)
	}
	userGet := rg.Group("/user")
	// TODO url encoded middleware
	{
		userGet.GET("/verify-email", userEmailVerifyHandl)
	}
	survey := rg.Group("/survey") // TODO example endpoints
	{
		survey.POST("/submit", surveySubmitHandl)
		survey.POST("/update", surveyUpdateHandl)
		survey.POST("/get", surveyGetHandl)
		survey.POST("/get-all", surveyGetAllHandl)
	}
	token := rg.Group("/token")
	{
		token.GET("/renew", tokenRenewHandl)
	}
}
