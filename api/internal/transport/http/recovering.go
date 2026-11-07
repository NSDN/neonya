package http

import (
	"encoding/json"
	"log/slog"
	stdhttp "net/http"
	"runtime/debug"
)

// internalErrorDocument 是 panic 兜底响应的 JSON 体，字段与 Huma 的错误体对齐。
type internalErrorDocument struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

// recovering 把 panic 变成 500 + problem+json，并把记录写进 slog（JSON 到 stdout）。
//
// 不用 chi 的 middleware.Recoverer：它把美化后的堆栈直接写到 stderr，
// 采集器只读 stdout 时会丢掉整条 panic 记录，且响应体是空的。
func recovering(logger *slog.Logger) func(stdhttp.Handler) stdhttp.Handler {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}
				// http.ErrAbortHandler 表示主动放弃响应，不该记成 500。
				if recovered == stdhttp.ErrAbortHandler {
					panic(recovered)
				}

				logger.Error("http handler panic",
					"request_id", requestIDFrom(request.Context()),
					"method", request.Method,
					"path", request.URL.Path,
					"panic", recovered,
					"stack", string(debug.Stack()),
				)

				writer.Header().Set("Content-Type", "application/problem+json")
				writer.WriteHeader(stdhttp.StatusInternalServerError)
				_ = json.NewEncoder(writer).Encode(internalErrorDocument{
					Title:  "Internal Server Error",
					Status: stdhttp.StatusInternalServerError,
					Detail: "internal error",
				})
			}()

			next.ServeHTTP(writer, request)
		})
	}
}
