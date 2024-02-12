package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"p2p/initializers"
	"p2p/models"
	"p2p/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

// SignUpUser
// RegisterUser
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Verificar se já existe um usuário com o email ou telefone
	existingUser := models.User{}

	if payload.Telephone != "" {
		err := ac.DB.Where("telephone = ?", payload.Telephone).First(&existingUser).Error
		if err == nil {
			// User already exists with the provided telephone
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email/telephone already exists"})
			return
		}
	}

	if payload.Email != "" {
		err := ac.DB.Where("email = ?", payload.Email).First(&existingUser).Error
		if err == nil {
			// Usuário já existe com o email fornecido
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email/telephone already exists"})
			return
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Handle other errors if needed
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Internal Server Error"})
			return
		}
		// If the error is ErrRecordNotFound, do nothing as it indicates the record wasn't found.
	}

	// Hash da senha, configuração do usuário e criação no banco de dados com Draft = true
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Agora, só depois de garantir que não há usuário existente, criamos o novo usuário
	now := time.Now()
	newUser := models.User{
		UserName: utils.GenerateRandomName(9),
		FullName: nil,
		Telephone: func() *string {
			if payload.Telephone != "" {
				return &payload.Telephone
			}
			return nil
		}(),
		Email: func() *string {
			if payload.Email != "" {
				email := strings.ToUpper(payload.Email)
				return &email
			}
			return nil
		}(),
		Password:            hashedPassword,
		Draft:               true,
		IsIdentityVerified:  false,
		IsEmailVerified:     false,
		IsTelephoneVerified: false,
		IsAddressVerified:   false,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	// Generate and Save OTP
	otp, err := ac.GenerateAndSaveOTP(payload.Email, payload.Telephone)
	if err != nil {
		// If there's an error generating OTP, remove the newly created user
		ac.DB.Delete(&newUser)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to generate OTP"})
		return
	}

	// Send Email with OTP
	emailData := utils.EmailData{
		URL:       otp,
		FirstName: "rosad_tests", // Define o primeiro nome do usuário aqui, se necessário
		Subject:   "Your OTP for account verification 2P2 AOA",
	}

	utils.SendEmail(&newUser, &emailData)

	message := "We sent an email with an OTP to " + *newUser.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

func (ac *AuthController) GenerateAndSaveOTP(email, telephone string) (string, error) {
	otp := utils.GenerateOTP()
	fmt.Println("Email-telephone", email, telephone)

	newOTPCode := models.OTPCode{
		Code:      otp,
		Email:     &email,
		Telephone: &telephone,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Expira em 1 dia
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := ac.DB.Create(&newOTPCode).Error
	if err != nil {
		return "", err
	}

	return otp, nil
}

func (ac *AuthController) ConfirmOTP(ctx *gin.Context) {
	var payload *models.OTPConfirmationInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Procurar o código OTP na tabela OTPCode
	var otpCode models.OTPCode
	result := ac.DB.Where("code = ? AND (email = ? OR telephone = ?) AND verified = false AND expires_at > ?", payload.Code, payload.Email, payload.Telephone, time.Now()).First(&otpCode)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid OTP code or expired"})
		return
	}

	// Verificar se o código tem exatamente 6 dígitos
	if len(payload.Code) != 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid OTP code format"})
		return
	}

	// Marcar o código OTP como verificado
	otpCode.Verified = true
	otpCode.UpdatedAt = time.Now()
	ac.DB.Save(&otpCode)

	// Atualizar o estado de Draft para false
	var user models.User
	err := ac.DB.Where("email = ? OR telephone = ?", strings.ToUpper(payload.Email), payload.Telephone).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "User not found"})
		return
	}

	user.Draft = false
	ac.DB.Save(&user)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OTP confirmed. User account activated"})
}

func (ac *AuthController) RequestNewOTP(ctx *gin.Context) {
	var payload *models.RequestNewOTPInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Verificar se o código anterior expirou ou não foi verificado
	var otpCode models.OTPCode
	result := ac.DB.Where("(email = ? OR telephone = ?) AND verified = false AND expires_at <= ?", payload.Email, payload.Telephone, time.Now()).Order("expires_at DESC").First(&otpCode)

	if result.Error == nil {
		// O código anterior expirou ou não foi verificado, podemos criar um novo
		otp := utils.GenerateOTP()

		// Salvar o novo código OTP na tabela OTPCode
		newOTPCode := models.OTPCode{
			Code:      otp,
			Email:     &payload.Email,
			Telephone: &payload.Telephone,
			Verified:  false,
			ExpiresAt: time.Now().Add(24 * time.Hour), // Expira em 1 dia
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		ac.DB.Create(&newOTPCode)

		// Enviar o novo OTP por e-mail ou SMS
		if payload.Email != "" {
			utils.SendOTPByEmail(payload.Email, otp)
		} else if payload.Telephone != "" {
			// Substitua com a função para enviar OTP por SMS
			// utils.SendOTPBySMS(payload.Telephone, otp, twilioConfig)
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "New OTP sent"})
		return
	}

	// Verificar se o usuário com o email ou telefone fornecido existe
	var existingUser models.User
	err := ac.DB.Where("email = ? OR telephone = ?", strings.ToUpper(payload.Email), payload.Telephone).First(&existingUser).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "User not found"})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Cannot request new OTP at the moment"})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	config, _ := initializers.LoadConfig(".")

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

// Refresh Access Token
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := initializers.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AuthController) VerifyEmail(ctx *gin.Context) {
	var payload models.OTPConfirmationInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var otpCode models.OTPCode
	if err := ac.DB.Where("code = ?", payload.Code).First(&otpCode).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exist"})
		return
	}

	if otpCode.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User already verified"})
		return
	}

	otpCode.Verified = true
	if err := ac.DB.Save(&otpCode).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to verify email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
}
