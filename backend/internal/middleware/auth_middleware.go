package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"Agora-BBS/internal/pkg/jwt"
	"Agora-BBS/internal/pkg/response"
)

// JWTAuth 拦截未带 Token 或 Token 失效的请求
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, 40101, "authorization header is required")
			c.Abort()
			return
		}

		// 按 'Bearer <token>' 格式解析
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, 40102, "authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		// 校验并解析 Token
		claims, err := jwt.ParseToken(parts[1], secret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 40103, "invalid or expired token")
			c.Abort()
			return
		}

		// 将解析出的身份上下文写入 Gin 上下文，供后续 Handler 使用
		c.Set("current_user_id", claims.UserID)
		c.Set("current_user_role", claims.Role)
		c.Next()
	}
}
