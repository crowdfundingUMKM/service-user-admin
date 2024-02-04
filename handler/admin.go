package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"service-user-admin/auth"
	"service-user-admin/core"
	"service-user-admin/database"
	"service-user-admin/helper"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

type userAdminHandler struct {
	userService core.Service
	authService auth.Service
}

var (
	storageClient *storage.Client
)

func NewUserHandler(userService core.Service, authService auth.Service) *userAdminHandler {
	return &userAdminHandler{userService, authService}
}

// Super Admin
func (h *userAdminHandler) GetLogtoAdmin(c *gin.Context) {
	// get data from middleware
	currentAdmin := c.MustGet("currentAdmin").(core.User)

	// id := os.Getenv("ADMIN_ID")
	if currentAdmin.RefAdmin == "MASTER" {
		content, err := os.ReadFile("./tmp/gin.log")
		if err != nil {
			response := helper.APIResponse("Failed to get log", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		// download with browser
		if c.Query("download") == "true" {
			c.Header("Content-Disposition", "attachment; filename=gin.log")
			c.Data(http.StatusOK, "application/octet-stream", content)
			return
		}

		c.String(http.StatusOK, string(content))
		return
	} else {
		response := helper.APIResponse("Your not Root Admin & Wrong Uri, cannot Access", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
}

// Get status service
func (h *userAdminHandler) ServiceHealth(c *gin.Context) {
	// check env open or not
	currentAdmin := c.MustGet("currentAdmin").(core.User)

	errEnv := godotenv.Load()
	if errEnv != nil {
		response := helper.APIResponse("Failed to get env for service investor", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// id := os.Getenv("ADMIN_ID")
	// if currentAdmin.RefAdmin == "MASTER" {
	// 	response := helper.APIResponse("Your not Admin, cannot Access1", http.StatusUnprocessableEntity, "error", nil)
	// 	c.JSON(http.StatusNotFound, response)
	// 	return
	// }
	errService := c.Errors
	if errService != nil {
		response := helper.APIResponse("Service Admin is not running", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	// check env open or not
	if currentAdmin.RefAdmin == "MASTER" {
		envVars := []string{
			"ADMIN_ID",
			"DB_USER",
			"DB_PASS",
			"DB_NAME",
			"DB_PORT",
			"INSTANCE_HOST",
			"SERVICE_HOST",
			"SERVICE_PORT",
			"JWT_SECRET",
			"STATUS_ACCOUNT",
		}

		data := make(map[string]interface{})
		for _, key := range envVars {
			data[key] = os.Getenv(key)
		}
		response := helper.APIResponse("Service Admin is running", http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
	} else {
		response := helper.APIResponse("Your not Admin MASTER, cannot Access", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
}

// Deactive admin user
func (h *userAdminHandler) DeactiveUser(c *gin.Context) {
	var input core.DeactiveUserInput
	currentAdmin := c.MustGet("currentAdmin").(core.User)

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("User Not Found", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// check id admin
	// id := os.Getenv("ADMIN_ID")
	if currentAdmin.RefAdmin == "MASTER" {
		// get id user

		// deactive user
		deactive, err := h.userService.DeactivateAccountUser(input, currentAdmin.UnixID)

		data := gin.H{
			"success_deactive": deactive,
		}

		if err != nil {
			response := helper.APIResponse("Failed to deactive user", http.StatusBadRequest, "error", data)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response := helper.APIResponse("User has been deactive", http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
	} else {
		response := helper.APIResponse("Your not Admin, cannot Access", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
}

func (h *userAdminHandler) ActiveUser(c *gin.Context) {
	var input core.DeactiveUserInput

	currentAdmin := c.MustGet("currentAdmin").(core.User)
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("User Not Found", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// check id admin
	// id := os.Getenv("ADMIN_ID")
	if currentAdmin.RefAdmin == "MASTER" {
		// get id user

		// deactive user
		active, err := h.userService.ActivateAccountUser(input, currentAdmin.UnixID)

		data := gin.H{
			"success_deactive": active,
		}

		if err != nil {
			response := helper.APIResponse("Failed to active user", http.StatusBadRequest, "error", data)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response := helper.APIResponse("User has been active", http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
	} else {
		response := helper.APIResponse("Your not Admin, cannot Access", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
}

// get all user by admin
func (h *userAdminHandler) GetAllUserData(c *gin.Context) {
	currentAdmin := c.MustGet("currentAdmin").(core.User)
	// id := os.Getenv("ADMIN_ID")
	if currentAdmin.RefAdmin == "MASTER" {
		users, err := h.userService.GetAllUsers()
		if err != nil {
			response := helper.APIResponse("Failed to get admin", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response := helper.APIResponse("List of user admin", http.StatusOK, "success", users)
		c.JSON(http.StatusOK, response)
	} else {
		response := helper.APIResponse("Your not Admin MASTER, cannot Access", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
}

// update user by admin
// func (h *userAdminHandler) UpdateUserByAdmin(c *gin.Context) {
// 	var input core.UpdateUserInput

// 	currentAdmin := c.MustGet("currentAdmin").(core.User)
// 	err := c.ShouldBindJSON(&input)
// 	if err != nil {
// 		errors := helper.FormatValidationError(err)
// 		errorMessage := gin.H{"errors": errors}

// 		response := helper.APIResponse("User Not Found", http.StatusUnprocessableEntity, "error", errorMessage)
// 		c.JSON(http.StatusUnprocessableEntity, response)
// 		return
// 	}
// 	// check id admin
// 	// id := os.Getenv("ADMIN_ID")
// 	if currentAdmin.RefAdmin == "MASTER" {
// 		// get id user by body unix_id target
// 		unixId := c.Param("unix_id")

// 		// deactive user
// 		update, err := h.userService.UpdateUserByUnixID(unixId, input)

// 		data := gin.H{
// 			"success_update": update,
// 		}

// 		if err != nil {
// 			dataError := gin.H{
// 				"errors": err.Error(),
// 			}
// 			response := helper.APIResponse("Failed to update user", http.StatusBadRequest, "error", dataError)
// 			c.JSON(http.StatusBadRequest, response)
// 			return
// 		}
// 		response := helper.APIResponse("User has been update", http.StatusOK, "success", data)
// 		c.JSON(http.StatusOK, response)
// 	} else {
// 		response := helper.APIResponse("Your not Admin, cannot Access", http.StatusUnprocessableEntity, "error", nil)
// 		c.JSON(http.StatusNotFound, response)
// 		return
// 	}
// }

// update password user by admin
func (h *userAdminHandler) UpdatePasswordByAdmin(c *gin.Context) {
	var inputID core.GetUserIdInput

	// check id is valid or not
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Update password failed", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData core.UpdatePasswordByAdminInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Update password failed, input data failure", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentAdmin").(core.User)

	updatedUser, err := h.userService.UpdatePasswordByAdmin(inputID.UnixID, inputData, currentUser.UnixID)
	if err != nil {
		response := helper.APIResponse("Update password failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := core.FormatterUserDetail(currentUser, updatedUser)
	response := helper.APIResponse("Password has been updated By Admin Master", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
	return
}

// delete user
func (h *userAdminHandler) DeleteUser(c *gin.Context) {
	var input core.DeleteUserInput

	currentAdmin := c.MustGet("currentAdmin").(core.User)
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("User Not Found", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// check id admin
	// id := os.Getenv("ADMIN_ID")
	if currentAdmin.RefAdmin == "MASTER" {
		// get id user

		// deactive user
		_, err := h.userService.DeleteUsers(input.UnixID)

		data := gin.H{
			"success_delete": true,
		}

		if err != nil {
			dataError := gin.H{
				"error": "User Not Found or Already Deleted",
			}
			response := helper.APIResponse("Failed to delete user", http.StatusBadRequest, "error", dataError)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response := helper.APIResponse("User has been delete", http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
	} else {
		response := helper.APIResponse("Your not Admin, cannot Access", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
}

// cen acces to veriefy token VerifyToken
func (h *userAdminHandler) VerifyToken(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(core.User)

	// check f account deactive
	if currentUser.StatusAccount == "deactive" {
		errorMessage := gin.H{"errors": "Your account is deactive by admin"}
		response := helper.APIResponse("Your token failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// if you logout you can't get user
	if currentUser.Token == "" {
		errorMessage := gin.H{"errors": "Your account is logout"}
		response := helper.APIResponse("Your token failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"success":  "Your token is valid",
		"admin_id": currentUser.UnixID,
	}

	response := helper.APIResponse("Successfuly get user by middleware", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

// User Admin
// Register User Admin
func (h *userAdminHandler) RegisterUser(c *gin.Context) {
	// tangkap input dari user
	// map input dari user ke struct RegisterUserInput
	// struct di atas kita passing sebagai parameter service

	var input core.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// generate token
	token, err := h.authService.GenerateToken(newUser.UnixID)
	if err != nil {
		if err != nil {
			response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	formatter := core.FormatterUser(newUser, token)

	if formatter.StatusAccount == "active" {
		_, err = h.userService.SaveToken(newUser.UnixID, token)

		if err != nil {
			response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := helper.APIResponse("Account has been registered and active", http.StatusOK, "success", formatter)
		c.JSON(http.StatusOK, response)
		return
	}

	data := gin.H{
		"status": "Account has been registered, but you must wait admin to active your account",
	}

	response := helper.APIResponse("Account has been registered but you must wait admin or review to active your account", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}

// Login User Admin
func (h *userAdminHandler) Login(c *gin.Context) {

	var input core.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// generate token
	token, err := h.authService.GenerateToken(loggedinUser.UnixID)

	// save toke to database
	_, err = h.userService.SaveToken(loggedinUser.UnixID, token)

	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// end save token to database

	if err != nil {
		if err != nil {
			response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	// check role acvtive and not send massage your account deactive
	if loggedinUser.StatusAccount == "deactive" {
		errorMessage := gin.H{"errors": "Your account is deactive by admin"}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	formatter := core.FormatterUser(loggedinUser, token)

	response := helper.APIResponse("Succesfuly loggedin", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

// Check Email Availability
func (h *userAdminHandler) CheckEmailAvailability(c *gin.Context) {
	var input core.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

// Check Phone Availability
func (h *userAdminHandler) CheckPhoneAvailability(c *gin.Context) {
	var input core.CheckPhoneInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Phone checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isPhoneAvailable, err := h.userService.IsPhoneAvailable(input)

	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.APIResponse("Phone checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isPhoneAvailable,
	}

	metaMessage := "Phone has been registered"

	if isPhoneAvailable {
		metaMessage = "Phone is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userAdminHandler) ServiceStart(c *gin.Context) {
	// check env open or not
	serviceStatus := "Service is active"

	errService := c.Errors
	if errService != nil {
		response := helper.APIResponse("Service investor is not running", http.StatusInternalServerError, "error", serviceStatus)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := helper.APIResponse("Service investor is running", http.StatusOK, "success", serviceStatus)
	c.JSON(http.StatusOK, response)
}

func (h *userAdminHandler) ServiceCheckDB(c *gin.Context) {
	// check env open or not
	// Check database status
	db := database.NewConnectionDB()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer sqlDB.Close()

	errDB := sqlDB.Ping()
	if errDB != nil {
		response := helper.APIResponse("Database is not running", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	currentUser := c.MustGet("currentUser").(core.User)

	data := gin.H{"user": "Wellcome to Service: " + currentUser.Name, "status": "Database is running"}

	response := helper.APIResponse("Service is running", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

// get user by middleware
func (h *userAdminHandler) GetUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(core.User)

	// check f account deactive
	if currentUser.StatusAccount == "deactive" {
		errorMessage := gin.H{"errors": "Your account is deactive by admin"}
		response := helper.APIResponse("Get user failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// if you logout you can't get user
	if currentUser.Token == "" {
		errorMessage := gin.H{"errors": "Your account is logout"}
		response := helper.APIResponse("Get user failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := core.FormatterUser(currentUser, "")

	response := helper.APIResponse("Successfuly get user by middleware", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

// update user by unix id
func (h *userAdminHandler) UpdateUser(c *gin.Context) {

	var inputData core.UpdateUserInput

	err := c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Update user failed, input data failure", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(core.User)

	// if you logout you can't get user
	if currentUser.Token == "" {
		errorMessage := gin.H{"errors": "Your account is logout"}
		response := helper.APIResponse("Get user failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updatedUser, err := h.userService.UpdateUserByUnixID(currentUser.UnixID, inputData)
	if err != nil {
		response := helper.APIResponse("Update user failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := core.FormatterUserDetail(currentUser, updatedUser)

	response := helper.APIResponse("User has been updated", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
	return
}

// update password by unix id
func (h *userAdminHandler) UpdatePassword(c *gin.Context) {

	var inputData core.UpdatePasswordInput

	err := c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Update password failed, input data failure", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(core.User)

	// if you logout you can't get user
	if currentUser.Token == "" {
		errorMessage := gin.H{"errors": "Your account is logout"}
		response := helper.APIResponse("Get user failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updatedUser, err := h.userService.UpdatePasswordByUnixID(currentUser.UnixID, inputData)
	if err != nil {
		response := helper.APIResponse("Update password failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// remove token in database
	_, err = h.userService.DeleteToken(currentUser.UnixID)
	if err != nil {
		response := helper.APIResponse("Update password failed & failed remove tokens", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := core.FormatterUserDetail(currentUser, updatedUser)
	response := helper.APIResponse("Password has been updated", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
	return
}

// get info id admin not use middleware
func (h *userAdminHandler) GetInfoAdminID(c *gin.Context) {
	var inputID core.GetUserIdInput

	// check id is valid or not
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed get user admin and status", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := h.userService.GetUserByUnixID(inputID.UnixID)
	if err != nil {
		response := helper.APIResponse("Get user failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := core.FormatterUserAdminID(user)

	response := helper.APIResponse("Successfuly get user id and status", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}

// Upload image
func (h *userAdminHandler) UploadAvatar(c *gin.Context) {
	f, _, err := c.Request.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(core.User)
	userID := currentUser.UnixID
	userName := currentUser.Name

	// if you logout you can't get user
	if currentUser.Token == "" {
		errorMessage := gin.H{"errors": "Your account is logout"}
		response := helper.APIResponse("Get user failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// initiate cloud storage os.Getenv("GCS_BUCKET")
	bucket := fmt.Sprintf("%s", os.Getenv("GCS_BUCKET"))
	subfolder := fmt.Sprintf("%s", os.Getenv("GCS_SUBFOLDER"))
	// var err error
	ctx := appengine.NewContext(c.Request)

	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("secret-keys.json"))

	if err != nil {
		// data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image to GCP", http.StatusBadRequest, "error", err)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	defer f.Close()

	objectName := fmt.Sprintf("%s/avatar-%s-%s", subfolder, userID, userName)
	sw := storageClient.Bucket(bucket).Object(objectName).NewWriter(ctx)

	if _, err := io.Copy(sw, f); err != nil {
		// data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image to GCP", http.StatusBadRequest, "error", err)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := sw.Close(); err != nil {
		// data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image to GCP", http.StatusBadRequest, "error", err)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	u, err := url.Parse("/" + bucket + "/" + sw.Attrs().Name)
	if err != nil {
		// data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image to GCP", http.StatusBadRequest, "error", err)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	path := u.String()

	// save avatar to database
	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfuly uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}

// Logout user
func (h *userAdminHandler) LogoutUser(c *gin.Context) {
	// get data from middleware
	currentUser := c.MustGet("currentUser").(core.User)

	// check if token is empty
	if currentUser.Token == "" {
		response := helper.APIResponse("Logout failed, your logout right now", http.StatusForbidden, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// delete token in database
	_, err := h.userService.DeleteToken(currentUser.UnixID)
	if err != nil {
		response := helper.APIResponse("Logout failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Logout success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
	return
}
