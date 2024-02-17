package routes

import (
	"p2p/controllers"
	"p2p/middleware"

	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("v1/auth")

	router.POST("/register", rc.authController.SignUpUser)
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/logout", middleware.DeserializeUser(), rc.authController.LogoutUser)
	router.POST("/register/otp", rc.authController.ConfirmOTP)
	router.POST("/register/resend/otp", rc.authController.RequestNewOTP)
	router.POST("/register/change/email_tephone", rc.authController.ChangeEmailOrTelephoneBeforeToBeVerified)

	router.GET("/verifyemail/:verificationCode", rc.authController.VerifyEmail)

}
