package middleware

import (
	"net/http"
	"os"
	"service-user-admin/auth"
	"service-user-admin/core"
	"service-user-admin/helper"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware(authService auth.Service, userService core.Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized MASTER", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized MASTER", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		// check if token is valid
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized ", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userUnixID := claim["unix_id"].(string)

		user, err := userService.GetUserByUnixID(userUnixID)
		if err != nil {
			response := helper.APIResponse("Unauthorized MASTER", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// check if user is admin
		if userUnixID != os.Getenv("ADMIN_ID") {
			data := gin.H{"status": "you are not MASTER admin"}
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", data)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("currentAdmin", user)
		c.Next()
	}
}
