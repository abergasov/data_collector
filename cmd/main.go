package main

import (
	"data_collector/pkg/config"
	"data_collector/pkg/logger"
	"data_collector/pkg/repository"
	"data_collector/pkg/routes"
	"data_collector/pkg/storage"
	"log"

	"github.com/valyala/fasthttp"

	"go.uber.org/zap"
)

var (
	appName   = "data_collector"
	buildTime = "_dev"
	buildHash = "_dev"
	confFile  = "common.yml"
)

func main() {
	err := logger.NewLogger()
	if err != nil {
		log.Fatal("error log init", err)
	}
	appConfig := config.InitConf(confFile)
	dbConnect := storage.InitDBConnect(appConfig)
	collectorM := repository.NewCollector(dbConnect)
	collectorSM := repository.NewCollectorSW(dbConnect)
	collectorSNG := repository.NewCollectorSNG(dbConnect)
	router := routes.InitRouter(appConfig, collectorSNG, collectorSM, collectorM, appName, buildHash, buildTime)
	logger.Info(
		"Server running on port",
		zap.String("port", appConfig.AppPort),
		zap.String("url", "http://localhost:"+appConfig.AppPort),
	)
	r := router.InitRoutes()
	err = fasthttp.ListenAndServe(":"+appConfig.AppPort, r.Handler)
	if err != nil {
		logger.Fatal("Router error", err)
	}
}
