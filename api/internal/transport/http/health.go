package http

import (
	"context"
	"log/slog"
	stdhttp "net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type healthDependencies struct {
	database DatabasePinger
	logger   *slog.Logger
	version  string
}

// healthResponse 是 /health 的响应体。
type healthResponse struct {
	State      string    `json:"state" doc:"ok 表示服务与数据库都可用"`
	Database   string    `json:"database" doc:"up / down"`
	Version    string    `json:"version"`
	ServerTime time.Time `json:"serverTime"`
}

type healthOutput struct {
	Status int
	Body   healthResponse
}

func registerHealthOperation(api huma.API, dependencies healthDependencies) {
	huma.Register(api, huma.Operation{
		OperationID:   "getHealth",
		Method:        stdhttp.MethodGet,
		Path:          "/health",
		Summary:       "服务健康检查",
		Description:   "数据库不可达时返回 503 与 degraded 状态。",
		Tags:          []string{"system"},
		DefaultStatus: stdhttp.StatusOK,
		Responses: map[string]*huma.Response{
			"503": {Description: "数据库不可达，响应体状态为 degraded"},
		},
	}, func(ctx context.Context, _ *struct{}) (*healthOutput, error) {
		serverTime := time.Now().UTC()

		healthy := false
		if dependencies.database != nil {
			err := dependencies.database.Ping(ctx)
			healthy = err == nil
			if err != nil {
				dependencies.logger.Error("health check failed", "error", err)
			}
		}

		if healthy {
			return &healthOutput{
				Status: stdhttp.StatusOK,
				Body: healthResponse{
					State:      "ok",
					Database:   "up",
					Version:    dependencies.version,
					ServerTime: serverTime,
				},
			}, nil
		}

		return &healthOutput{
			Status: stdhttp.StatusServiceUnavailable,
			Body: healthResponse{
				State:      "degraded",
				Database:   "down",
				Version:    dependencies.version,
				ServerTime: serverTime,
			},
		}, nil
	})
}
