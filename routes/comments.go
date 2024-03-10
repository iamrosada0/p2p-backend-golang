package routes

import (
	"p2p/controllers"
	"p2p/middleware"

	"github.com/gin-gonic/gin"
)

type CommentsRouteController struct {
	CommentsController controllers.CommentsController
}

func NewRouteCommentsController(CommentsController controllers.CommentsController) CommentsRouteController {
	return CommentsRouteController{CommentsController}
}

func (uc *CommentsRouteController) CommentsRoute(rg *gin.RouterGroup) {
	router := rg.Group("comments")
	// Rotas relacionadas aos comentários do usuário
	router.POST("/comments", middleware.DeserializeUser(), uc.CommentsController.CreateComments)
	router.GET("/comments", middleware.DeserializeUser(), uc.CommentsController.GetCommentss)
	router.DELETE("/comments/:id", middleware.DeserializeUser(), uc.CommentsController.DeleteComments)

}
