package router

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/fx"
	"kits/api/src/common/configs"
	appvalidator "kits/api/src/common/validator"
	httpserver "kits/api/src/present/http"
	"kits/api/src/present/http/controllers"
	"kits/api/src/present/http/middlewares"
	"net/http"
)

type RoutersIn struct {
	fx.In
	Config     *configs.Config
	HttpServer *httpserver.HttpServer
	AuthCtrl   *controllers.AuthCtrl
	UserCtrl   *controllers.UserCtrl
}

func registerMiddlewares(in RoutersIn) {
	root := in.HttpServer.Root
	cf := in.Config

	// Recovery
	root.Use(gin.CustomRecovery(middlewares.Recovery))
	root.Use(middlewares.NewTrackHandler())
	root.Use(otelgin.Middleware(cf.Server.Name))
	root.Use(middlewares.NewLogRequest())
}

func RegisterRouters(in RoutersIn) {
	appvalidator.RegisterGinValidator()
	registerMiddlewares(in)

	root := in.HttpServer.Root
	v1 := root.Group("v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		registerAuthRouters(v1, in)
	}
}

func registerAuthRouters(root *gin.RouterGroup, in RoutersIn) {
	authRouter := root.Group("auth")
	{
		authRouter.POST("/sign-up", in.UserCtrl.Create)
		authRouter.POST("/sign-in", in.AuthCtrl.Login)
		authRouter.POST("/refresh-token", in.AuthCtrl.RefreshToken)
	}
}
