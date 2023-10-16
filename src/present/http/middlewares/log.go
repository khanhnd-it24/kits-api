package middlewares

import (
	"github.com/gin-gonic/gin"
	"kits/api/src/common/logger"
)

func NewLogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		logReq := map[string]interface{}{
			"http.path":       path,
			"http.method":     c.Request.Method,
			"http.query":      c.Request.URL.RawQuery,
			"http.user_agent": c.Request.UserAgent(),
		}

		logCtx := logger.WithContextValue(c.Request.Context(), logReq)
		c.Request = c.Request.WithContext(logCtx)
		c.Next()
	}
}
