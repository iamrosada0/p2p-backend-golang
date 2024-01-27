package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"p2p/api"

	"p2p/internal/user/entity"
	"p2p/internal/user/infra/repository"
	"p2p/internal/user/usecase"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/lib/pq"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	dbPath := "./db/main.db"
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			panic(err)
		}

		file, err := os.Create(dbPath)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	// Create Gorm connection
	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = gormDB.AutoMigrate(&entity.User{})
	if err != nil {
		panic(err)
	}
	// Create repositories and use cases
	userRepository := repository.NewUserRepositoryPostgres(gormDB)
	createUserUsecase := usecase.NewCreateUserUseCase(userRepository)
	listUsersUsecase := usecase.NewGetAllUsersUseCase(userRepository)
	deleteUserUsecase := usecase.NewDeleteUserUseCase(userRepository)
	getUserByIDUsecase := usecase.NewGetUserByIDUseCase(userRepository)
	updateUserUsecase := usecase.NewUpdateUserUseCase(userRepository)

	// Create handlers
	userHandlers := api.NewUserHandlers(createUserUsecase, listUsersUsecase, deleteUserUsecase, getUserByIDUsecase, updateUserUsecase)

	// Set up Gin router
	router := gin.Default()

	// Set up user routes
	userHandlers.SetupRoutes(router)

	// Start the server
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println(err)
	}
}
