package controllers

import (
	"errors"
	"net/http"
	"p2p/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FriendController struct {
	DB *gorm.DB
}

func NewFriendController(DB *gorm.DB) FriendController {
	return FriendController{DB}
}

func (fc *FriendController) AddFriend(ctx *gin.Context) {
	// AddFriend é um controlador para lidar com a adição de um amigo.
	// Este método recebe uma solicitação HTTP contendo o número de telefone do amigo a ser adicionado
	// e verifica se o usuário autenticado tem permissão para adicionar amigos.
	//Eu e a pessoa que estou adicionar precisamos ter o nosso numero gravado

	// Obtém o ID do usuário autenticado que está realizando a ação
	userID := ctx.MustGet("userID").(uint)

	// Obtém o número de telefone do amigo a ser adicionado dos parâmetros da solicitação
	friendPhone := ctx.Param("friendPhone")

	// Verifica se o número de telefone do amigo não está vazio
	if friendPhone == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Friend's phone number is required"})
		return
	}

	// Verifica se o usuário atual existe
	currentUser := models.User{}
	if err := fc.DB.First(&currentUser, "id = ?", userID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user"})
		return
	}

	// Verifica se o amigo existe com base no número de telefone fornecido
	friendUser := models.User{}
	if err := fc.DB.First(&friendUser, "telephone = ?", friendPhone).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
		return
	}

	// Verifica se ambos os usuários têm números de telefone válidos
	if currentUser.Telephone == nil || friendUser.Telephone == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Both users must have phone numbers"})
		return
	}

	// Verifica se os números de telefone correspondem
	if *currentUser.Telephone != friendPhone {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone numbers don't match"})
		return
	}

	// Verifica se já existe uma relação de amizade entre os usuários
	var existingFriendship models.Friend
	result := fc.DB.First(&existingFriendship, "user_id = ? AND friend_id = ?", userID, friendUser.ID)
	if result.Error == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Friendship already exists"})
		return
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing friendship"})
		return
	}

	// Cria uma nova entrada na tabela de amigos
	newFriendship := models.Friend{
		UserID:    userID,
		FriendID:  friendUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := fc.DB.Create(&newFriendship).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add friend"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Friend added successfully"})
}

func (fc *FriendController) RemoveFriend(ctx *gin.Context) {
	// RemoveFriend é um controlador para lidar com a remoção de um amigo.
	// Este método recebe uma solicitação HTTP contendo o ID do amigo a ser removido
	// e verifica se o usuário autenticado tem permissão para remover amigos.

	// Obtém o ID do usuário autenticado que está realizando a ação
	userID := ctx.MustGet("userID").(uint)

	// Obtém o ID do amigo a ser removido dos parâmetros da solicitação
	friendID := ctx.Param("friendID")

	// Verifica se o ID do amigo não está vazio
	if friendID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Friend ID is required"})
		return
	}

	// Verifica se o ID do amigo é um número válido
	friendIDUint, err := strconv.ParseUint(friendID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	// Verifica se a amizade existe e pertence ao usuário autenticado
	var friendship models.Friend
	result := fc.DB.First(&friendship, "user_id = ? AND friend_id = ?", userID, friendIDUint)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Friendship not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check friendship"})
		}
		return
	}

	// Remove a amizade do banco de dados
	if err := fc.DB.Delete(&friendship).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove friend"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Friend removed successfully"})
	//Quando terminar de remover como amigo, deve ser enviado uma mensagem para a outra pessoa dizem que
	//ele foi removido pelo (xyz)/ entao irei adicionar rabbitmq
}

func (fc *FriendController) GetFriends(ctx *gin.Context) {
	// GetFriends é um controlador para obter a lista de amigos de um usuário.
	// Este método recebe uma solicitação HTTP e retorna a lista de amigos do usuário autenticado.

	// Obtém o ID do usuário autenticado que está realizando a ação
	userID := ctx.MustGet("userID").(uint)

	// Variável para armazenar a lista de amigos do usuário
	var friends []models.User

	// Consulta ao banco de dados para obter a lista de amigos do usuário
	// Utilizamos uma subquery para encontrar todos os IDs de amigos do usuário
	// e, em seguida, buscamos os detalhes desses usuários
	if err := fc.DB.Raw("SELECT * FROM users WHERE id IN (SELECT friend_id FROM friends WHERE user_id = ?)", userID).Scan(&friends).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get friends"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": friends})
}
