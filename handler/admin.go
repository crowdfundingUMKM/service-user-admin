package handler

import (
	"io/ioutil"
	"net/http"
	"os"
	"service-user-admin/admin"
	"service-user-admin/auth"
	"service-user-admin/helper"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type userAdminHandler struct {
	userService admin.Service
	authService auth.Service
}

func NewUserHandler(userService admin.Service, authService auth.Service) *userAdminHandler {
	return &userAdminHandler{userService, authService}
}

func (h *userAdminHandler) GetLogtoAdmin(c *gin.Context) {

	id := os.Getenv("ADMIN_ID")
	if c.Param("id") == id {
		content, err := ioutil.ReadFile("./log/gin.log")
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
	} else {
		response := helper.APIResponse("Your not Root Admin, cannot Access", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
}

// for admin get env
func (h *userAdminHandler) ServiceHealth(c *gin.Context) {
	// check env open or not
	errEnv := godotenv.Load()
	if errEnv != nil {
		response := helper.APIResponse("Failed to get env for service investor", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	id := os.Getenv("ADMIN_ID")
	if c.Param("id") != id {
		response := helper.APIResponse("Your not Admin, cannot Access", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
	// check env open or not
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	instance_host := os.Getenv("INSTANCE_HOST")
	service_host := os.Getenv("SERVICE_HOST")
	service_port := os.Getenv("SERVICE_PORT")
	jwt_secret := os.Getenv("JWT_SECRET")
	status_account := os.Getenv("STATUS_ACCOUNT")
	admin_id := os.Getenv("ADMIN_ID")

	data := map[string]interface{}{
		"db_user":          db_user,
		"db_pass":          db_pass,
		"db_name":          db_name,
		"db_port":          db_port,
		"db_instance_host": instance_host,
		"service_host":     service_host,
		"service_port":     service_port,
		"jwt_secret":       jwt_secret,
		"status_account":   status_account,
		"admin_id":         admin_id,
	}
	errService := c.Errors
	if errService != nil {
		response := helper.APIResponse("Service Admin is not running", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := helper.APIResponse("Service Admin is running", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userAdminHandler) RegisterUser(c *gin.Context) {
	// tangkap input dari user
	// map input dari user ke struct RegisterUserInput
	// struct di atas kita passing sebagai parameter service

	var input admin.RegisterUserInput

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

	formatter := admin.FormatterUser(newUser, token)

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

func (h *userAdminHandler) Login(c *gin.Context) {

	var input admin.LoginInput

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
	formatter := admin.FormatterUser(loggedinUser, token)

	response := helper.APIResponse("Succesfuly loggedin", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userAdminHandler) CheckEmailAvailability(c *gin.Context) {
	var input admin.CheckEmailInput

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

func (h *userAdminHandler) CheckPhoneAvailability(c *gin.Context) {
	var input admin.CheckPhoneInput

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

// update user by unix id
func (h *userAdminHandler) UpdateUser(c *gin.Context) {
	var inputID admin.GetUserDetailInput

	// check id is valid or not
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Update user failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData admin.UpdateUserInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Update user failed, input data failure", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(admin.User)

	if currentUser.UnixID != inputID.UnixID {
		response := helper.APIResponse("Update user failed, because you are not auth", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	updatedUser, err := h.userService.UpdateUserByUnixID(currentUser.UnixID, inputData)
	if err != nil {
		response := helper.APIResponse("Update user failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := admin.FormatterUserDetail(currentUser, updatedUser)

	response := helper.APIResponse("User has been updated", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
	return
}
