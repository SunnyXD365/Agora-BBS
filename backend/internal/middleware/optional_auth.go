package middleware

import (
	"strings"

	"agora-backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func OptionalJWT(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			if claims, err := jwt.ParseToken(parts[1], jwtSecret); err == nil {
				c.Set("userID", claims.UserID)
			}
		}
		c.Next()
	}
}
