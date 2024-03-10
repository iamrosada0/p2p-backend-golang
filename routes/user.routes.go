package routes

import (
	"p2p/controllers"
	"p2p/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("users")

	// Rotas relacionadas ao perfil do usuário
	router.GET("/me", middleware.DeserializeUser(), uc.userController.GetMe)
	router.POST("/profile", middleware.DeserializeUser(), uc.userController.CreateProfile)
	router.PUT("/profile/:id", middleware.DeserializeUser(), uc.userController.UpdateProfile)
	// Adicione mais rotas conforme necessário

}
