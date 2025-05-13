package middleware

import (
	"kondangin-backend/internal/service"
	"kondangin-backend/internal/utils"
	"log"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWTToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		username := claims.Subject
		log.Printf("username middleware : %s", username)

		// 🔍 Panggil UserService untuk ambil user dari DB
		user, err := userService.GetUserByUsername(username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Set user ke context biar bisa dipakai di handler
		c.Set("user", user)

		c.Next()
	}
}
