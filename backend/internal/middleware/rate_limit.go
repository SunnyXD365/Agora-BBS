package middleware

import (
	"fmt"
	"net/http"
	"time"

	"agora-backend/internal/cache"
	"agora-backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func RateLimit(store *cache.Store, scope string, limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate:%s:%s", scope, c.ClientIP())
		allowed, err := store.Allow(c.Request.Context(), key, limit, window)
		if err == nil && !allowed {
			response.Error(c, http.StatusTooManyRequests, 42901, "too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}
