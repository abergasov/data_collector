package routes

import (
	"data_collector/pkg/config"
	"data_collector/pkg/middleware"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type AppRouter struct {
	appBuildHash string
	appBuildTime string
	appName      string
	EchoEngine   *echo.Echo
	config       *config.AppConfig
	collector    ICollector
}

func InitRouter(cnf *config.AppConfig, c ICollector, appName, appHash, appBuild string) *AppRouter {
	router := &AppRouter{
		EchoEngine:   echo.New(),
		config:       cnf,
		appName:      appName,
		appBuildHash: appHash,
		appBuildTime: appBuild,
		collector:    c,
	}
	router.EchoEngine.Use(middleware.ZapLogger())
	return router
}

func (ar *AppRouter) InitRoutes() *echo.Echo {
	ar.EchoEngine.GET("/", ar.Ping)
	ar.EchoEngine.GET("/collect", ar.Handler)
	return ar.EchoEngine
}

// Handler
func (ar *AppRouter) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		OK        bool      `json:"ok"`
		Now       time.Time `json:"now"`
		BuildHash string    `json:"build_hash"`
		BuildTime string    `json:"build_time"`
		AppName   string    `json:"app_name"`
	}{
		OK:        true,
		Now:       time.Now(),
		AppName:   ar.appName,
		BuildHash: ar.appBuildHash,
		BuildTime: ar.appBuildTime,
	})
}

func (ar *AppRouter) Handler(c echo.Context) error {
	ar.collector.HandleEvent(1, "awd")
	return c.JSON(http.StatusOK, `{"ok":true}`)
}
