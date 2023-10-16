package bootstrap

import (
	"go.uber.org/fx"
	"kits/api/src/core/services"
)

func BuildServices() fx.Option {
	return fx.Options(
		fx.Provide(services.NewUserService),
		fx.Provide(services.NewAuthService),
	)
}
