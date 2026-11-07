package http

import (
	"context"
	"crypto/rand"
	stdhttp "net/http"
)

const requestIDHeaderName = "X-Request-Id"

type requestIDContextKey struct{}

// requestID 给每个请求生成新的 ID 并写进响应头。
//
// 不读入站的同名头：chi 的 middleware.RequestID 会沿用客户端传来的值，
// 外部随手伪造一个就能让多条请求共用同一条日志线索，审计与排障都不可信。
func requestID(next stdhttp.Handler) stdhttp.Handler {
	return stdhttp.HandlerFunc(func(writer stdhttp.ResponseWriter, request *stdhttp.Request) {
		identifier := rand.Text()
		writer.Header().Set(requestIDHeaderName, identifier)

		requestContext := context.WithValue(request.Context(), requestIDContextKey{}, identifier)
		next.ServeHTTP(writer, request.WithContext(requestContext))
	})
}

// requestIDFrom 取出当前请求的 ID，没设置过则返回空串。
func requestIDFrom(ctx context.Context) string {
	identifier, _ := ctx.Value(requestIDContextKey{}).(string)
	return identifier
}
