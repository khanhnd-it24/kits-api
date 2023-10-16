package bootstrap

import (
	"context"
	"go.uber.org/fx"
	"kits/api/src/common/logger"
	httpserver "kits/api/src/present/http"
	"kits/api/src/present/http/router"
)

func BuildHttpServer() fx.Option {
	return fx.Options(
		fx.Provide(httpserver.NewHttpServer),
		fx.Invoke(func(lc fx.Lifecycle, s *httpserver.HttpServer) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := s.Start(ctx)
					if err != nil {
						logger.Error(ctx, err, "http server start failed")
						return err
					}
					logger.Info(ctx, "http server started")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					err := s.Stop(ctx)
					if err != nil {
						logger.Error(ctx, err, "http server stop failed")
						return err
					}
					logger.Info(ctx, "http server stopped")
					return nil

				},
			})
		}),
		fx.Invoke(router.RegisterRouters),
	)
}
