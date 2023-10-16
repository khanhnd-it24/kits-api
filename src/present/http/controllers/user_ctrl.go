package controllers

import (
	"github.com/gin-gonic/gin"
	"kits/api/src/common/fault"
	"kits/api/src/common/logger"
	"kits/api/src/core/domains"
	"kits/api/src/core/services"
	"kits/api/src/present/http/requests"
	httputil "kits/api/src/present/http/util"
)

type UserCtrl struct {
	userService *services.UserService
}

func NewUserCtrl(userService *services.UserService) *UserCtrl {
	return &UserCtrl{userService: userService}
}

func (ctrl *UserCtrl) Create(c *gin.Context) {
	ctxReq := c.Request.Context()
	caller := "UserCtrl.Create"

	var req requests.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		wErr := fault.ConvertValidatorErr(err)

		logger.Warn(ctxReq, wErr, "[%v] invalid params %+v", caller, req)
		httputil.ServeErrResponse(c, wErr)
		return
	}

	insertedId, err := ctrl.userService.Create(ctxReq, &domains.UserCreate{
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		logger.Warn(ctxReq, err, "[%v] failed to create user", caller)
		httputil.ServeErrResponse(c, err)
		return
	}

	res := map[string]int64{"id": insertedId}
	httputil.ServeSuccessResponse(c, res)
}
