package main

import (
	"flag"
	"go.uber.org/fx"
	"kits/api/src/bootstrap"
)

func createApp(pathCf string) *fx.App {
	return fx.New(
		bootstrap.BuildFxLogger(),
		bootstrap.BuildAppConfig(pathCf),
		bootstrap.BuildLogger(),
		bootstrap.BuildCrypto(),
		bootstrap.BuildDatabases(),
		bootstrap.BuildServices(),
		bootstrap.BuildControllers(),
		bootstrap.BuildHttpServer(),
	)
}

func main() {
	var pathConfig string

	flag.StringVar(&pathConfig, "config", "configs/config.yaml", "path to config file")
	flag.Parse()

	app := createApp(pathConfig)
	// when call run, app automatic handle graceful shutdown
	app.Run()
}
