package main

import (
	"fmt"
	"log"
	"os"
	"service-user-admin/admin"
	"service-user-admin/auth"
	"service-user-admin/database"
	"service-user-admin/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// setup repository
	db := database.NewConnectionDB()
	userAdminRepository := admin.NewRepository(db)

	// SETUP SERVICE
	userAdminService := admin.NewService(userAdminRepository)
	authService := auth.NewService()

	// setup handler
	userHandler := handler.NewUserHandler(userAdminService, authService)

	// END SETUP

	// RUN SERVICE
	router := gin.Default()
	api := router.Group("api/v1")

	// Rounting admin-health

	// Rounting admin
	api.POST("register_admin", userHandler.RegisterUser)

	// end Rounting
	url := fmt.Sprintf("%s:%s", os.Getenv("SERVICE_HOST"), os.Getenv("SERVICE_PORT"))
	router.Run(url)

}
