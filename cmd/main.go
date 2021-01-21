package main

import (
	"data_collector/pkg/config"
	"data_collector/pkg/logger"
	"data_collector/pkg/routes"
	"go.uber.org/zap"
	"log"
)

func main() {
	appConfig := config.InitConf()
	err := logger.NewLogger()
	if err != nil {
		log.Fatal("error log init", err)
	}
	router := routes.InitRouter(appConfig)
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
