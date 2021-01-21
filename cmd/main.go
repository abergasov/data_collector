package main

import (
	"data_collector/pkg/config"
	"data_collector/pkg/logger"
	"data_collector/pkg/repository"
	"data_collector/pkg/routes"
	"log"

	"github.com/valyala/fasthttp"

	"go.uber.org/zap"
)

var (
	appName   = "data_collector"
	buildTime = "_dev"
	buildHash = "_dev"
)

func main() {
	err := logger.NewLogger()
	if err != nil {
		log.Fatal("error log init", err)
	}
	appConfig := config.InitConf("common.yml")
	collector := repository.NewCollector()
	router := routes.InitRouter(appConfig, collector, appName, buildHash, buildTime)
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
