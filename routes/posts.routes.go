package routes

import (
	"p2p/controllers"
	"p2p/middleware"

	"github.com/gin-gonic/gin"
)

type PostsRouteController struct {
	PostsController controllers.PostsController
}

func NewRoutePostsController(PostsController controllers.PostsController) PostsRouteController {
	return PostsRouteController{PostsController}
}

func (uc *PostsRouteController) PostsRoute(rg *gin.RouterGroup) {
	router := rg.Group("Postss")

	// Rotas relacionadas às postagens do usuário
	router.GET("/posts", middleware.DeserializeUser(), uc.PostsController.GetPosts)
	router.POST("/posts", middleware.DeserializeUser(), uc.PostsController.CreatePost)
	router.DELETE("/posts/:id", middleware.DeserializeUser(), uc.PostsController.DeletePost)
	// Adicione mais rotas conforme necessário

}
