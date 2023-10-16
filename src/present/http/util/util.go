package httputil

import (
	"github.com/gin-gonic/gin"
	"kits/api/src/common/fault"
	"kits/api/src/common/track"
	"net/http"
)

const (
	SuccessCode    = 0
	SuccessKey     = "Success"
	SuccessMessage = "success"
)

type Response struct {
	Code    int64       `json:"code"`
	Key     string      `json:"key"`
	Body    interface{} `json:"body"`
	Message string      `json:"message"`
	TrackId string      `json:"track_id"`
}

func ServeSuccessResponse(c *gin.Context, body interface{}) {
	trackId := track.GetTrackId(c.Request.Context())

	c.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Key:     SuccessKey,
		Body:    body,
		Message: SuccessMessage,
		TrackId: trackId,
	})
}

func ServeErrResponse(c *gin.Context, err error, statusCodes ...int) {
	var statusCode int
	if len(statusCodes) > 0 {
		statusCode = statusCodes[0]
	} else {
		statusCode = fault.GetHttpStatusCode(err)
	}

	trackId := track.GetTrackId(c.Request.Context())

	errRes := Response{
		Key:     fault.GetKey(err),
		Code:    fault.GetCode(err),
		Message: fault.GetDescription(err),
		TrackId: trackId,
		Body:    nil,
	}

	c.JSON(statusCode, errRes)
}
