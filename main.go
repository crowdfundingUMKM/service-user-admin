package main

import (
	"fmt"
	"log"
	"os"
	"service-user-admin/admin"
	"service-user-admin/auth"
	"service-user-admin/database"
	"service-user-admin/handler"
	"service-user-admin/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// setup log
	// L.InitLog()

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
	// config cors
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowWildcard = true
	// hasil permintaan CORS akan disimpan di cache browser selama 12 jam sebelum browser melakukan permintaan baru ke server.
	config.MaxAge = 12 * time.Hour
	config.AllowFiles = true
	router.Use(cors.New(config))

	// group api
	api := router.Group("api/v1")

	// Rounting admin-health Root Admin
	// api.GET("/log_service_admin/:id", userHandler.GetLogtoAdmin)
	api.GET("/service_status/:id", userHandler.ServiceHealth)

	// Rounting admin
	api.POST("/email_check", userHandler.CheckEmailAvailability)
	api.POST("/phone_check", userHandler.CheckPhoneAvailability)
	api.POST("/register_admin", userHandler.RegisterUser)
	api.POST("/login_admin", userHandler.Login)

	// get user by unix_id
	api.GET("/get_user", middleware.AuthMiddleware(authService, userAdminService), userHandler.GetUser)
	api.PUT("/update_admin/:unix_id", middleware.AuthMiddleware(authService, userAdminService), userHandler.UpdateUser)

	api.POST("/logout_admin", middleware.AuthMiddleware(authService, userAdminService), userHandler.Logout)

	// end Rounting
	url := fmt.Sprintf("%s:%s", os.Getenv("SERVICE_HOST"), os.Getenv("SERVICE_PORT"))
	router.Run(url)

}
