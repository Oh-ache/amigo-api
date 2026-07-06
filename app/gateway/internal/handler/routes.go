package handler

import (
	"net/http"

	"amigo-api/app/gateway/internal/proxy"
	"amigo-api/app/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	manager := proxy.NewManager(serverCtx.Config.Routes)

	// 健康检查
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/health",
		Handler: manager.HealthHandler(),
	})

	// API 代理：覆盖 1-5 层路径深度，6 个 HTTP 方法
	depths := []string{"/api/:a", "/api/:a/:b", "/api/:a/:b/:c", "/api/:a/:b/:c/:d", "/api/:a/:b/:c/:d/:e", "/api/:a/:b/:c/:d/:e/:f"}
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}

	for _, depth := range depths {
		routes := make([]rest.Route, 0, len(methods))
		for _, method := range methods {
			routes = append(routes, rest.Route{Method: method, Path: depth, Handler: manager.ServeHTTP})
		}
		server.AddRoutes(routes)
	}
}
