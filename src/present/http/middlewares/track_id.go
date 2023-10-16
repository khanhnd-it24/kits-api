package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kits/api/src/common/logger"
)

func NewTrackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var trackId string

		requestId := c.GetHeader("X-REQUEST-ID")
		if requestId != "" {
			trackId = requestId
		} else {
			trackId = uuid.New().String()
		}

		ctxReq := c.Request.Context()
		ctxWithTrackId := context.WithValue(ctxReq, "track_id", trackId)
		loggerCtx := logger.WithContextValue(ctxWithTrackId, map[string]interface{}{"track_id": trackId})

		c.Request = c.Request.WithContext(loggerCtx)
		c.Next()
	}
}
