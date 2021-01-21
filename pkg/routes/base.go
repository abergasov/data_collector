package routes

import (
	"data_collector/pkg/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type AppRouter struct {
	EchoEngine *echo.Echo
	config     *config.AppConfig
}

func InitRouter(cnf *config.AppConfig) *AppRouter {
	router := &AppRouter{
		EchoEngine: echo.New(),
		config:     cnf,
	}
	return router
}

func (ar *AppRouter) InitRoutes() *echo.Echo {
	ar.EchoEngine.GET("/", ar.Ping)
	ar.EchoEngine.GET("/collect", ar.Ping)
	return ar.EchoEngine
}

// Handler
func (ar *AppRouter) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		OK  bool      `json:"ok"`
		Now time.Time `json:"now"`
	}{
		OK:  true,
		Now: time.Now(),
	})
}
