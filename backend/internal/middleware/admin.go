package middleware

import (
	"database/sql"
	"net/http"

	"agora-backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func AdminOnly(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		verified, _ := c.Get("adminVerified")
		if verified != true {
			response.Error(c, http.StatusForbidden, 40391, "administrator email verification required")
			c.Abort()
			return
		}
		var role, status string
		if err := db.QueryRowContext(c.Request.Context(), `SELECT role,status FROM users WHERE id=$1`, c.GetInt64("userID")).Scan(&role, &status); err != nil || role != "admin" || status != "active" {
			response.Error(c, http.StatusForbidden, 40390, "administrator permission required")
			c.Abort()
			return
		}
		c.Next()
	}
}
