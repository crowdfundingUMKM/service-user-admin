package main

import (
	"fmt"
	"log"
	"os"
	"service-user-admin/auth"
	"service-user-admin/config"
	"service-user-admin/core"
	"service-user-admin/database"
	"service-user-admin/handler"
	"service-user-admin/middleware"

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
	// config.InitLog()

	// setup repository
	db := database.NewConnectionDB()
	userAdminRepository := core.NewRepository(db)

	// SETUP SERVICE
	userAdminService := core.NewService(userAdminRepository)
	authService := auth.NewService()

	// setup handler
	userHandler := handler.NewUserHandler(userAdminService, authService)

	// END SETUP

	// RUN SERVICE
	router := gin.Default()

	// setup cors
	corsConfig := config.InitCors()
	router.Use(cors.New(corsConfig))

	// group api
	api := router.Group("api/v1")

	// Rounting admin-health Root Admin
	api.GET("/log_service_admin/:admin_id", middleware.AdminMiddleware(authService, userAdminService), userHandler.GetLogtoAdmin)
	api.GET("/service_status/:admin_id", middleware.AdminMiddleware(authService, userAdminService), userHandler.ServiceHealth)
	api.POST("/deactive_user/:admin_id", middleware.AdminMiddleware(authService, userAdminService), userHandler.DeactiveUser)
	api.POST("/active_user/:admin_id", middleware.AdminMiddleware(authService, userAdminService), userHandler.ActiveUser)
	api.DELETE("/delete_user/:admin_id", middleware.AdminMiddleware(authService, userAdminService), userHandler.DeleteUser)
	api.PUT("/update_user_by_admin/:admin_id/:unix_id", middleware.AdminMiddleware(authService, userAdminService), userHandler.UpdateUserByAdmin)
	api.PUT("/update_password_by_admin/:admin_id/:unix_id", middleware.AdminMiddleware(authService, userAdminService), userHandler.UpdatePasswordByAdmin)

	api.GET("/get_all_user_by_admin", middleware.AdminMiddleware(authService, userAdminService), userHandler.GetAllUserData)

	// make endoint to change statusbyadmin MASTER

	// make endpoint SoftDeletebyadmin MASTER

	// update sql if admin master change must add time update

	// can access to get prove token
	// verify token
	api.GET("/verifyTokenAdmin", middleware.AuthMiddleware(authService, userAdminService), userHandler.VerifyToken)

	// Rounting admin
	api.POST("/email_check", userHandler.CheckEmailAvailability)
	api.POST("/phone_check", userHandler.CheckPhoneAvailability)
	api.POST("/register_admin", userHandler.RegisterUser)
	api.POST("/login_admin", userHandler.Login)

	//make service health for investor
	api.GET("/service_start", userHandler.ServiceStart)
	api.GET("/service_check", middleware.AuthMiddleware(authService, userAdminService), userHandler.ServiceCheckDB)

	// route give information to user about admin
	api.GET("/admin/getAdminID/:unix_id", userHandler.GetInfoAdminID)

	// get user by unix_id
	api.GET("/get_user", middleware.AuthMiddleware(authService, userAdminService), userHandler.GetUser)
	api.PUT("/update_profile", middleware.AuthMiddleware(authService, userAdminService), userHandler.UpdateUser)
	api.PUT("/update_password", middleware.AuthMiddleware(authService, userAdminService), userHandler.UpdatePassword)

	api.POST("/logout_admin", middleware.AuthMiddleware(authService, userAdminService), userHandler.LogoutUser)

	//make Upload image profile user token and Update image avatar
	api.POST("/upload_avatar", middleware.AuthMiddleware(authService, userAdminService), userHandler.UploadAvatar)

	//make update image profile user by unix_id

	//make delete image profile user by unix_id

	// Notifikasi user admin

	// Create notifikasi user admin with auth

	// Update notifikasi user admin with auth

	// Delete notifikasi user admin with auth

	// Update Status notifikasi user admin with auth

	// end Rounting
	url := fmt.Sprintf("%s:%s", os.Getenv("SERVICE_HOST"), os.Getenv("SERVICE_PORT"))
	router.Run(url)

}
