package controllers

import (
	"net/http"
	"p2p/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostsController struct {
	DB *gorm.DB
}

func NewPostController(DB *gorm.DB) PostsController {
	return PostsController{DB}
}

// AddPost é um controlador para adicionar uma nova postagem.
func (pc *PostsController) CreatePost(ctx *gin.Context) {
	// Obter os dados da postagem do corpo da solicitação
	var postData models.Post
	if err := ctx.ShouldBindJSON(&postData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post data"})
		return
	}

	// Adicionar detalhes adicionais à postagem, como a hora de criação
	postData.CreatedAt = time.Now()

	// Persistir a postagem no banco de dados
	if err := pc.DB.Create(&postData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add post"})
		return
	}

	// Responder com sucesso
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": postData})
}

// DeletePost é um controlador para excluir uma postagem existente.
func (pc *PostsController) DeletePost(ctx *gin.Context) {
	// Obter o ID da postagem a ser excluída dos parâmetros da solicitação
	postIDStr := ctx.Param("id")

	// Converter o ID da postagem para um tipo uint
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Verificar se a postagem existe no banco de dados
	var post models.Post
	result := pc.DB.First(&post, postID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Verificar se o usuário autenticado é o proprietário da postagem
	userID := ctx.MustGet("userID").(uint)
	if post.UserID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Excluir a postagem do banco de dados
	if err := pc.DB.Delete(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	// Responder com sucesso
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Post deleted successfully"})
}

// GetPosts é um controlador para obter todas as postagens.
func (pc *PostsController) GetPosts(ctx *gin.Context) {
	// Recuperar todas as postagens do banco de dados
	var posts []models.Post
	if err := pc.DB.Find(&posts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	// Responder com todas as postagens recuperadas
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": posts})
}
