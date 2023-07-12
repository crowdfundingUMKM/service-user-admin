package handler

import (
	"io/ioutil"
	"net/http"
	"os"
	"service-user-admin/auth"
	"service-user-admin/core"
	"service-user-admin/helper"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type userAdminHandler struct {
	userService core.Service
	authService auth.Service
}

func NewUserHandler(userService core.Service, authService auth.Service) *userAdminHandler {
	return &userAdminHandler{userService, authService}
}

// Super Admin
func (h *userAdminHandler) GetLogtoAdmin(c *gin.Context) {
	// get data from middleware
	currentAdmin := c.MustGet("currentAdmin").(core.User)

	id := os.Getenv("ADMIN_ID")
	if c.Param("admin_id") == currentAdmin.UnixID && c.Param("admin_id") == id {
		content, err := ioutil.ReadFile("./tmp/gin.log")
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

	id := os.Getenv("ADMIN_ID")
	if c.Param("admin_id") != id {
		response := helper.APIResponse("Your not Admin, cannot Access", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
	errService := c.Errors
	if errService != nil {
		response := helper.APIResponse("Service Admin is not running", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	// check env open or not
	if c.Param("admin_id") == currentAdmin.UnixID && c.Param("admin_id") == id {
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
		response := helper.APIResponse("Service Admin is running", http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
	} else {
		response := helper.APIResponse("Your not Admin, cannot Access", http.StatusUnprocessableEntity, "error", nil)
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
	id := os.Getenv("ADMIN_ID")
	if c.Param("admin_id") == currentAdmin.UnixID && c.Param("admin_id") == id {
		// get id user

		// deactive user
		deactive, err := h.userService.DeactivateAccountUser(input)

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
	id := os.Getenv("ADMIN_ID")
	if c.Param("admin_id") == currentAdmin.UnixID && c.Param("admin_id") == id {
		// get id user

		// deactive user
		active, err := h.userService.ActivateAccountUser(input)

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

// update user by admin
func (h *userAdminHandler) UpdateUserByAdmin(c *gin.Context) {
	var input core.UpdateUserInput

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
	id := os.Getenv("ADMIN_ID")
	if c.Param("admin_id") == currentAdmin.UnixID && c.Param("admin_id") == id {
		// get id user by body unix_id target
		unixId := c.Param("unix_id")

		// deactive user
		update, err := h.userService.UpdateUserByUnixID(unixId, input)

		data := gin.H{
			"success_update": update,
		}

		if err != nil {
			dataError := gin.H{
				"errors": err.Error(),
			}
			response := helper.APIResponse("Failed to update user", http.StatusBadRequest, "error", dataError)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response := helper.APIResponse("User has been update", http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
	} else {
		response := helper.APIResponse("Your not Admin, cannot Access", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

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
	id := os.Getenv("ADMIN_ID")
	if c.Param("admin_id") == currentAdmin.UnixID && c.Param("admin_id") == id {
		// get id user

		// deactive user
		delete, err := h.userService.DeleteUsers(input.UnixID)

		data := gin.H{
			"success_delete": delete,
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

	formatter := core.FormatterUser(currentUser, "")

	response := helper.APIResponse("Successfuly get user by middleware", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

// update user by unix id
func (h *userAdminHandler) UpdateUser(c *gin.Context) {
	var inputID core.GetUserIdInput

	// check id is valid or not
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Update user failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData core.UpdateUserInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Update user failed, input data failure", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(core.User)

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

	formatter := core.FormatterUserDetail(currentUser, updatedUser)

	response := helper.APIResponse("User has been updated", http.StatusOK, "success", formatter)
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

// Logout user
func (h *userAdminHandler) Logout(c *gin.Context) {
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
