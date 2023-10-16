package bootstrap

import (
	"go.uber.org/fx"
	"kits/api/src/common/configs"
	hashprovider "kits/api/src/common/crypto/hash"
	"kits/api/src/common/logger"
)

func BuildAppConfig(pathConfig string) fx.Option {
	return fx.Provide(func() (*configs.Config, error) {
		return configs.NewAppConfig(pathConfig)
	})
}

func BuildLogger() fx.Option {
	return fx.Options(
		fx.Provide(logger.NewLogger),
		fx.Invoke(func(l logger.Logger) {
			logger.SetLogger(l)
		}),
	)
}

func BuildCrypto() fx.Option {
	return fx.Options(
		fx.Provide(hashprovider.NewHashProvider),
	)
}
