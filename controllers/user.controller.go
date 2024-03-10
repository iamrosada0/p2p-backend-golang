package controllers

import (
	"net/http"
	"strconv"
	"time"

	"p2p/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		ID:        currentUser.ID,
		UserName:  currentUser.UserName,
		Telephone: *currentUser.Telephone,
		Email:     *currentUser.Email,
		Biography: *currentUser.Biography,
		Photo:     *currentUser.Photo,

		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	var user models.User
	if err := uc.DB.First(&user, "id = ?", userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
}
func (uc *UserController) CreateProfile(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user profile"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": user})
}

func (uc *UserController) UpdateProfile(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := uc.DB.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updatedUser models.User
	if err := ctx.ShouldBindJSON(&updatedUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualize os campos relevantes do perfil de usuário
	user.FullName = updatedUser.FullName
	user.Email = updatedUser.Email
	user.Telephone = updatedUser.Telephone
	user.IsIdentityVerified = updatedUser.IsIdentityVerified
	user.IsEmailVerified = updatedUser.IsEmailVerified
	user.IsTelephoneVerified = updatedUser.IsTelephoneVerified
	user.IsAddressVerified = updatedUser.IsAddressVerified
	user.Draft = updatedUser.Draft

	user.Photo = updatedUser.Photo
	user.Biography = updatedUser.Biography
	user.UpdatedAt = time.Now()

	if err := uc.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
}
