// Package http 是 HTTP 传输层：只做编解码、校验与错误映射，不判权限、不碰 SQL。
package http

import (
	"context"
	"log/slog"
	stdhttp "net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// APIBasePath 是所有接口的固定前缀；浏览器与服务端客户端的 baseUrl 都用它。
const APIBasePath = "/api/v1"

// ContractVersion 是 OpenAPI 文档里的接口版本，与路径前缀 /v1 对应。
const ContractVersion = "1.0.0"

// DatabasePinger 是健康检查需要的最小数据库能力。
type DatabasePinger interface {
	Ping(ctx context.Context) error
}

// RouterOptions 是组装路由所需的依赖。
type RouterOptions struct {
	Logger *slog.Logger
	// Database 允许为 nil：导出 OpenAPI 不需要连库。
	Database DatabasePinger
	Version  string
	// TrustedProxyCount 是本服务与客户端之间的代理跳数。
	// 0 表示没有代理，只认 TCP 连接来源，绝不解析 X-Forwarded-For。
	TrustedProxyCount int
}

// Router 同时暴露 HTTP handler 与 Huma API，后者用于导出契约。
type Router struct {
	Handler *chi.Mux
	API     huma.API
}

// NewRouter 组装中间件与全部 operation。
func NewRouter(options RouterOptions) Router {
	rootRouter := chi.NewMux()
	rootRouter.Use(
		requestID,
		// 必须排在 requestLogging 之前：日志要读 context 里的客户端 IP。
		clientIPMiddleware(options.TrustedProxyCount),
		requestLogging(options.Logger),
		recovering(options.Logger),
	)

	configuration := huma.DefaultConfig("neonya API", ContractVersion)
	// 路径保持相对（/health），前缀由 servers 声明，前端 baseUrl 直接用 /api/v1。
	configuration.OpenAPI.Servers = []*huma.Server{{URL: APIBasePath}}

	apiRouter := chi.NewMux()
	api := humachi.New(apiRouter, configuration)

	registerHealthOperation(api, healthDependencies{
		database: options.Database,
		logger:   options.Logger,
		version:  options.Version,
	})

	rootRouter.Mount(APIBasePath, apiRouter)

	return Router{Handler: rootRouter, API: api}
}

// clientIPMiddleware 只在明确知道前面有几层代理时才解析 X-Forwarded-For。
// 跳数配高了拿不到 IP（空值），配低了客户端就能伪造 IP。
func clientIPMiddleware(trustedProxyCount int) func(stdhttp.Handler) stdhttp.Handler {
	if trustedProxyCount > 0 {
		return middleware.ClientIPFromXFFTrustedProxies(trustedProxyCount)
	}
	return middleware.ClientIPFromRemoteAddr
}
