package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"agora-backend/internal/pkg/jwt"
	"agora-backend/internal/pkg/response"
)

func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, 40101, "authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, 40102, "authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1], jwtSecret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 40103, "invalid or expired token")
			c.Abort()
			return
		}

		// 将解析出的 userID 写入上下文，后端的 Handler 可以随时取出
		c.Set("userID", claims.UserID)
		c.Next()
	}
}