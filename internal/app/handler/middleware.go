package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
)

const prefix = "Bearer"

const (
	AdminRole = iota
	UserRole
	GuestRole
)

// RoleMiddleware проверяет роль пользователя и проверяет блеклист
func (h *Handler) RoleMiddleware(allowedRoles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractTokenFromHeader(c.Request)

		// Check if token exists (optional, depending on your guest access policy)
		if tokenString == "" {
			// Handle guest access if no token is present
			if !containsInt(allowedRoles, GuestRole) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			// Set user ID to 0 (or a placeholder value) for guests
			c.Set("user_id", float64(0))
			c.Set("is_admin", false)
			c.Next()
			return
		}

		// Parse and validate token (unchanged)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if h.Repository.IsBlackListed(userID, tokenString) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token in blacklist"})
			return
		}

		// Set user ID and admin flag from token claims (unchanged)
		c.Set("user_id", userID)
		isAdmin, ok := claims["is_admin"].(bool)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}
		c.Set("is_admin", isAdmin)

		// Check role against allowed roles (unchanged)
		if !isRoleAllowed(isAdmin, allowedRoles) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
				"msg":   "role not allowed",
			})
			return
		}

		// Middleware successful, proceed with handler
		c.Next()
	}
}

func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}
	if strings.Split(bearerToken, " ")[0] != prefix {
		return ""
	}
	return strings.Split(bearerToken, " ")[1]
}

func isRoleAllowed(isAdmin bool, allowedRoles []int) bool {
	for _, role := range allowedRoles {
		if isAdmin == (role == AdminRole) {
			return true
		}
	}
	return false

}

func containsInt(arr []int, val int) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

//
