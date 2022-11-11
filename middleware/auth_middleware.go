package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/letenk/use_deal_user/models/web"
	"github.com/letenk/use_deal_user/service"
)

// Function for auth middleware
func AuthMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get header with name `Authorization`
		authHeader := c.GetHeader("Authorization")

		// If inside authHeader doesn't have `Bearer`
		if !strings.Contains(authHeader, "Bearer") {
			// Create format response
			response := web.JSONResponseWithoutData(http.StatusUnauthorized, "error", "Unauthorized")
			// Stop process and return response
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// If there is, create new variable with empty string value
		tokenString := ""
		// Split authHeader with white space
		arrayToken := strings.Split(authHeader, " ")
		// If length arrayToken is same the 2
		if len(arrayToken) == 2 {
			// Get arrayToken with index 1 / only token jwt
			tokenString = arrayToken[1]
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, errors.New("invalid token")
			}

			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		// If error
		if err != nil {
			// Create format response
			response := web.JSONResponseWithoutData(http.StatusUnauthorized, "error", "Unauthorized")
			// Stop process and return response
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Get payload token
		claim, ok := token.Claims.(jwt.MapClaims)
		// If not `ok` and token invalid
		if !ok || !token.Valid {
			// Create format response
			response := web.JSONResponseWithoutData(http.StatusUnauthorized, "error", "Unauthorized")
			// Stop process and return response
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Get payload `user_id` and convert to `string`
		userId := claim["user_id"].(string)

		// Find user on db with service
		user, err := userService.GetOne(userId)
		// If error
		if err != nil {
			// Create format response
			response := web.JSONResponseWithoutData(http.StatusUnauthorized, "error", "Unauthorized")
			// Stop process and return response
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Set user to context with name `currentUser`
		c.Set("currentUser", user)
	}
}
