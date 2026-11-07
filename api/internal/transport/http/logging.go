package http

import (
	"log/slog"
	stdhttp "net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// requestLogging 输出 JSON 访问日志。
//
// 客户端 IP 与请求 ID 都由前面的中间件放进 context，这里只读：
// chi 的 RealIP 因为无条件信任 X-Forwarded-For 而废弃（可伪造）。
func requestLogging(logger *slog.Logger) func(stdhttp.Handler) stdhttp.Handler {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
			startedAt := time.Now()
			recorder := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

			next.ServeHTTP(recorder, request)

			logger.Info("http request",
				"request_id", requestIDFrom(request.Context()),
				"client_ip", middleware.GetClientIP(request.Context()),
				"method", request.Method,
				"path", request.URL.Path,
				"status", recorder.Status(),
				"bytes", recorder.BytesWritten(),
				"duration_ms", time.Since(startedAt).Milliseconds(),
			)
		})
	}
}
