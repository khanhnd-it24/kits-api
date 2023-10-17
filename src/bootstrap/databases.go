package bootstrap

import (
	"context"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"kits/api/src/common/cache"
	"kits/api/src/common/logger"
	apppostgres "kits/api/src/common/postgres"
	infracache "kits/api/src/infra/cache"
	"kits/api/src/infra/postgres/repos"
)

func BuildDatabases() fx.Option {
	return fx.Options(
		fx.Provide(apppostgres.NewPostgresProvider),
		fx.Invoke(func(lc fx.Lifecycle, dbProvider *apppostgres.DBProvider) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info(ctx, "db started")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					err := dbProvider.Stop(ctx)
					if err != nil {
						logger.Error(ctx, err, "db stop failed")
						return err
					}
					logger.Info(ctx, "db stopped")
					return nil
				},
			})
		}),
		fx.Provide(func(dbProvider *apppostgres.DBProvider) *gorm.DB {
			return dbProvider.DB()
		}),
		BuildRepos(),
		BuildCaches(),
	)
}

func BuildRepos() fx.Option {
	return fx.Options(
		fx.Provide(repos.NewUserRepo),
		fx.Provide(repos.NewRefreshTokenRepo),
	)
}

func BuildCaches() fx.Option {
	return fx.Options(
		fx.Provide(cache.NewRedisClient),
		fx.Provide(infracache.NewAccessTokenCache),
	)
}
