package controllers

import (
	"net/http"
	"p2p/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentsController struct {
	DB *gorm.DB
}

func NewCommentsController(DB *gorm.DB) CommentsController {
	return CommentsController{DB}
}

// CreateComments é um controlador para criar um novo comentário.
func (cc *CommentsController) CreateComments(ctx *gin.Context) {
	var comment models.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Realize qualquer validação adicional dos dados do comentário, se necessário

	if err := cc.DB.Create(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": comment})
}

// DeleteComments é um controlador para excluir um comentário existente.
func (cc *CommentsController) DeleteComments(ctx *gin.Context) {
	commentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var comment models.Comment
	if err := cc.DB.First(&comment, commentID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if err := cc.DB.Delete(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Comment deleted successfully"})
}

// GetCommentss é um controlador para obter todos os comentários.
func (cc *CommentsController) GetCommentss(ctx *gin.Context) {
	var comments []models.Comment
	if err := cc.DB.Find(&comments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": comments})
}
