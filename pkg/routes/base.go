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
	str500             = []byte(`{"ok":false}`)
)

type AppRouter struct {
	appBuildInfo   []byte
	FastHttpEngine *router.Router
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
		FastHttpEngine: router.New(),
		config:         cnf,
		collector:      c,
		appBuildInfo:   b,
	}
}

func (ar *AppRouter) InitRoutes() *router.Router {
	//ar.FastHttpEngine.GET("/", fasthttp.CompressHandler(ar.Index))
	ar.FastHttpEngine.GET("/", ar.Index)
	//ar.FastHttpEngine.POST("/collect", fasthttp.CompressHandler(ar.Collect))
	ar.FastHttpEngine.POST("/collect", ar.Collect)
	return ar.FastHttpEngine
}

func (ar *AppRouter) Index(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.SetStatusCode(http.StatusOK)
	ctx.Write(ar.appBuildInfo)
}

func (ar *AppRouter) Collect(ctx *fasthttp.RequestCtx) {
	event := &PayloadMessage{}
	err := event.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	if ar.collector.HandleEvent(event.ID, event.Label) {
		ctx.SetStatusCode(http.StatusOK)
		ctx.Write(strOK)
		return
	}
	ctx.SetStatusCode(http.StatusInternalServerError)
	ctx.Write(str500)
}
