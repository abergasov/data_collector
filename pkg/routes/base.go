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
	collectorSNG   ICollector // single map
	collectorSM    ICollector // syncMap
	collectorM     ICollector // multi map
}

func InitRouter(cnf *config.AppConfig, cSNG, cSM, cM ICollector, appName, appHash, appBuild string) *AppRouter {
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
		collectorSNG:   cSNG,
		collectorSM:    cSM,
		collectorM:     cM,
		appBuildInfo:   b,
	}
}

func (ar *AppRouter) InitRoutes() *router.Router {
	//ar.FastHTTPEngine.GET("/", fasthttp.CompressHandler(ar.Index))
	ar.FastHTTPEngine.GET("/", ar.Index)
	//ar.FastHTTPEngine.POST("/collect", fasthttp.CompressHandler(ar.Collect))
	ar.FastHTTPEngine.POST("/collect_single_map", ar.CollectSingleMap)
	ar.FastHTTPEngine.POST("/collect_multi_map", ar.CollectSingleMap)
	ar.FastHTTPEngine.POST("/collect_sync_map", ar.CollectSyncMap)
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
	stat := ar.collectorM.GetState()
	statBytes, _ := json.Marshal(stat)
	ctx.Write(statBytes)
}

func (ar *AppRouter) CollectSyncMap(ctx *fasthttp.RequestCtx) {
	event := &PayloadMessage{}
	err := event.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	ar.collectorSM.HandleEvent(event.ID, event.Label)
	ar.finishRequesto(ctx)
}

func (ar *AppRouter) CollectMultiMap(ctx *fasthttp.RequestCtx) {
	event := &PayloadMessage{}
	err := event.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	ar.collectorM.HandleEvent(event.ID, event.Label)
	ar.finishRequesto(ctx)
}

func (ar *AppRouter) CollectSingleMap(ctx *fasthttp.RequestCtx) {
	event := &PayloadMessage{}
	err := event.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	ar.collectorSNG.HandleEvent(event.ID, event.Label)
	ar.finishRequesto(ctx)
}

func (ar *AppRouter) finishRequesto(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(strOK)
}
