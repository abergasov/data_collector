package main

import (
	"data_collector/pkg/config"
	"data_collector/pkg/logger"
	"data_collector/pkg/repository"
	"data_collector/pkg/routes"
	"log"

	"go.uber.org/zap"
)

var (
	appName   = "data_collector"
	buildTime = "_dev"
	buildHash = "_dev"
)

func main() {
	err := logger.NewLogger()
	appConfig := config.InitConf("common.yml")
	if err != nil {
		log.Fatal("error log init", err)
	}
	collector := repository.NewCollector()
	router := routes.InitRouter(appConfig, collector, appName, buildHash, buildTime)
	// Start server
	logger.Info(
		"Server running on port",
		zap.String("port", appConfig.AppPort),
		zap.String("url", "http://localhost:"+appConfig.AppPort),
	)
	err = router.InitRoutes().Start(":" + appConfig.AppPort)
	if err != nil {
		logger.Fatal("Router error", err)
	}
}
