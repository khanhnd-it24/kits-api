package controllers

import (
	"github.com/gin-gonic/gin"
	"kits/api/src/common/fault"
	"kits/api/src/common/logger"
	"kits/api/src/core/domains"
	"kits/api/src/core/services"
	"kits/api/src/present/http/requests"
	"kits/api/src/present/http/responses"
	httputil "kits/api/src/present/http/util"
)

type AuthCtrl struct {
	authService *services.AuthService
}

func NewAuthCtrl(authService *services.AuthService) *AuthCtrl {
	return &AuthCtrl{authService: authService}
}

func (ctrl *AuthCtrl) Login(c *gin.Context) {
	ctxReq := c.Request.Context()
	caller := "AuthUserCtrl.Login"

	var req requests.AuthLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		wErr := fault.ConvertValidatorErr(err)
		logger.Warn(ctxReq, wErr, "[%v] invalid params %+v", caller, req)
		httputil.ServeErrResponse(c, wErr)
		return
	}

	token, user, err := ctrl.authService.Login(ctxReq, &domains.AuthUser{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		logger.Warn(ctxReq, err, "[%v] failed to create token", caller)
		httputil.ServeErrResponse(c, err)
		return
	}

	httputil.ServeSuccessResponse(c, responses.UserLoginRes{
		Token: &responses.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresIn:    token.ExpiresIn,
		},
		User: responses.UserFromDomain(user),
	})
}

func (ctrl *AuthCtrl) RefreshToken(c *gin.Context) {
	ctxReq := c.Request.Context()
	caller := "AuthUserCtrl.RefreshToken"

	var req requests.AuthRefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		wErr := fault.ConvertValidatorErr(err)
		logger.Warn(ctxReq, wErr, "[%v] invalid params %+v", caller, req)
		httputil.ServeErrResponse(c, wErr)
		return
	}

	token, err := ctrl.authService.RefreshToken(ctxReq, &domains.AuthRefreshToken{
		Token: req.RefreshToken,
	})

	if err != nil {
		logger.Warn(ctxReq, err, "[%v] failed to refresh token", caller)
		httputil.ServeErrResponse(c, err)
		return
	}

	httputil.ServeSuccessResponse(c, responses.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	})
}
