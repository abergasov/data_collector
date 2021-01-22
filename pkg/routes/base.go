package routes

import (
	"data_collector/pkg/config"
	"encoding/json"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
	strOK              = []byte(`{"ok":true}`)
)

type AppRouter struct {
	appBuildInfo   []byte
	FastHTTPEngine *router.Router
	config         *config.AppConfig
	collector      ICollector
}

func InitRouter(cnf *config.AppConfig, c ICollector, appName, appHash, appBuild string) *AppRouter {
	buildInfo := struct {
		OK        bool   `json:"ok"`
		BuildHash string `json:"build_hash"`
		BuildTime string `json:"build_time"`
		AppName   string `json:"app_name"`
	}{
		OK:        true,
		AppName:   appName,
		BuildHash: appHash,
		BuildTime: appBuild,
	}
	b, _ := json.Marshal(buildInfo)
	return &AppRouter{
		FastHTTPEngine: router.New(),
		config:         cnf,
		collector:      c,
		appBuildInfo:   b,
	}
}

func (ar *AppRouter) InitRoutes() *router.Router {
	//ar.FastHTTPEngine.GET("/", fasthttp.CompressHandler(ar.Index))
	ar.FastHTTPEngine.GET("/", ar.Index)
	//ar.FastHTTPEngine.POST("/collect", fasthttp.CompressHandler(ar.Collect))
	ar.FastHTTPEngine.POST("/collect", ar.Collect)
	ar.FastHTTPEngine.GET("/state", ar.State)
	return ar.FastHTTPEngine
}

func (ar *AppRouter) Index(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(ar.appBuildInfo)
}

func (ar *AppRouter) State(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.SetStatusCode(http.StatusOK)
	stat := ar.collector.GetState()
	statBytes, _ := json.Marshal(stat)
	ctx.Write(statBytes)
}

func (ar *AppRouter) Collect(ctx *fasthttp.RequestCtx) {
	event := &PayloadMessage{}
	err := event.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ar.collector.HandleEvent(event.ID, event.Label)
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(strOK)
}
