package bootstrap

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"kits/api/src/common/logger"
)

func BuildFxLogger() fx.Option {
	return fx.WithLogger(
		func() fxevent.Logger {
			return CreateFxEventLogger()
		},
	)
}

func CreateFxEventLogger() fxevent.Logger {
	l := logger.NewFxLogger()
	return &logWrapper{Logger: l}
}

type logWrapper struct {
	Logger *zap.SugaredLogger
}

func (zw *logWrapper) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
	case *fxevent.OnStartExecuted:
	case *fxevent.OnStopExecuting:
	case *fxevent.OnStopExecuted:
	case *fxevent.Supplied:
	case *fxevent.Provided:
		if e.Err != nil {
			scoped := make([]string, 0)
			scoped = append(scoped, e.OutputTypeNames...)
			scoped = append(scoped, "err", e.Err.Error())
			zw.Logger.Errorw("error encountered while applying options", scoped)
		}
	case *fxevent.Invoking:
	case *fxevent.Invoked:
	case *fxevent.Stopping:
	case *fxevent.Stopped:
		if e.Err != nil {
			zw.Logger.Errorw("stop failed", "err", e.Err)
		} else {
			zw.Logger.Info("service stopped")
		}
	case *fxevent.RollingBack:
	case *fxevent.RolledBack:
	case *fxevent.Started:
		if e.Err != nil {
			zw.Logger.Errorw("start failed", "err", e.Err)
		} else {
			zw.Logger.Info("service started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			zw.Logger.Errorw("custom logger initialization failed", "err", e.Err)
		}
	}
}
