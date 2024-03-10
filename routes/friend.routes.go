package routes

import (
	"p2p/controllers"
	"p2p/middleware"

	"github.com/gin-gonic/gin"
)

type FriendRouteController struct {
	FriendController controllers.FriendController
}

func NewRouteFriendController(FriendController controllers.FriendController) FriendRouteController {
	return FriendRouteController{FriendController}
}

func (uc *FriendRouteController) FriendRoute(rg *gin.RouterGroup) {
	router := rg.Group("friends")

	// Rotas relacionadas aos amigos do usuário
	router.GET("/friends", middleware.DeserializeUser(), uc.FriendController.GetFriends)
	router.POST("/friends", middleware.DeserializeUser(), uc.FriendController.AddFriend)
	router.DELETE("/friends/:id", middleware.DeserializeUser(), uc.FriendController.RemoveFriend)
	// Adicione mais rotas conforme necessário

	// Adicione mais rotas conforme necessário
}
