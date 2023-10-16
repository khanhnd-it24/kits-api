package bootstrap

import (
	"go.uber.org/fx"
	"kits/api/src/present/http/controllers"
)

func BuildControllers() fx.Option {
	return fx.Options(
		fx.Provide(controllers.NewUserCtrl),
		fx.Provide(controllers.NewAuthCtrl),
	)
}
