// Command api 是论坛的 HTTP 服务。子命令只有 serve（起服务）与 spec（导出 OpenAPI）。
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	stdhttp "net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	humayaml "github.com/danielgtaylor/huma/v2/yaml"

	"github.com/nsdn/neonya/api/internal/config"
	"github.com/nsdn/neonya/api/internal/platform/postgres"
	transporthttp "github.com/nsdn/neonya/api/internal/transport/http"
)

// version 由构建时的 -ldflags "-X main.version=..." 覆盖。
var version = "dev"

func main() {
	// 从命令的第二个参数取意思是忽略程序本身，也就是说是传输子命令。
	if err := run(os.Args[1:]); err != nil {
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).Error("api failed", "error", err)
		os.Exit(1)
	}
}

func run(arguments []string) error {
	if len(arguments) == 0 {
		return fmt.Errorf("usage: api <serve|spec>")
	}

	switch arguments[0] {
	case "serve":
		return serve()
	case "spec": // 用来写 `openapi.yaml`。See: `task gen`.
		return writeSpecification(arguments[1:])
	default:
		return fmt.Errorf("unknown command %q, want serve or spec", arguments[0])
	}
}

func serve() error {
	configuration, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	shutdownSignal, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()

	pool, err := postgres.NewPool(shutdownSignal, configuration.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	router := transporthttp.NewRouter(transporthttp.RouterOptions{
		Logger:            logger,
		Database:          pool,
		Version:           version,
		TrustedProxyCount: configuration.TrustedProxyCount,
	})

	server := &stdhttp.Server{
		Addr:              configuration.HTTPAddress,
		Handler:           router.Handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	logger.Info("api listening",
		"address", configuration.HTTPAddress,
		"environment", configuration.ApplicationEnvironment,
		"version", version,
	)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, stdhttp.ErrServerClosed) {
			return fmt.Errorf("serve: %w", err)
		}
	case <-shutdownSignal.Done():
	}

	logger.Info("api shutting down")

	shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelShutdown()

	return server.Shutdown(shutdownContext)
}

func writeSpecification(arguments []string) error {
	flags := flag.NewFlagSet("spec", flag.ContinueOnError)
	outputPath := flags.String("output", "", "写入文件路径；留空写标准输出")
	checkOnly := flags.Bool("check", false, "只校验 -output 指向的文件是否为最新契约，不写盘")
	if err := flags.Parse(arguments); err != nil {
		return err
	}

	if *checkOnly && *outputPath == "" {
		return fmt.Errorf("-check 需要配合 -output 指向要校验的文件")
	}

	document, err := renderSpecification()
	if err != nil {
		return err
	}

	if *checkOnly {
		return verifySpecification(*outputPath, document)
	}

	if *outputPath == "" {
		_, err := os.Stdout.Write(document)
		return err
	}

	if directory := filepath.Dir(*outputPath); directory != "." {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			return fmt.Errorf("create output directory: %w", err)
		}
	}

	return os.WriteFile(*outputPath, document, 0o644)
}

// renderSpecification 把当前路由注册导成 OpenAPI 3.1 的 YAML 字节。
func renderSpecification() ([]byte, error) {
	router := transporthttp.NewRouter(transporthttp.RouterOptions{
		Logger:  slog.New(slog.NewJSONHandler(io.Discard, nil)),
		Version: version,
	})

	jsonDocument, err := json.MarshalIndent(router.API.OpenAPI(), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal openapi json: %w", err)
	}

	var yamlDocument bytes.Buffer
	if err := humayaml.Convert(&yamlDocument, bytes.NewReader(jsonDocument)); err != nil {
		return nil, fmt.Errorf("convert openapi to yaml: %w", err)
	}

	return yamlDocument.Bytes(), nil
}

// verifySpecification 比对磁盘上的契约与当前代码导出的结果，只读不写。
func verifySpecification(path string, expected []byte) error {
	committed, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取 %s: %w（先跑 task gen 生成）", path, err)
	}

	// 只比内容：checkout 的换行转换（core.autocrlf）不该算成契约漂移。
	committedLines := strings.Split(strings.ReplaceAll(string(committed), "\r\n", "\n"), "\n")
	expectedLines := strings.Split(strings.ReplaceAll(string(expected), "\r\n", "\n"), "\n")

	for index := 0; index < len(committedLines) && index < len(expectedLines); index++ {
		if committedLines[index] != expectedLines[index] {
			return fmt.Errorf("%s 第 %d 行与当前代码不一致：磁盘 %q，导出 %q（跑 task gen 刷新）",
				path, index+1, committedLines[index], expectedLines[index])
		}
	}

	if len(committedLines) != len(expectedLines) {
		return fmt.Errorf("%s 行数不一致：磁盘 %d 行，导出 %d 行（跑 task gen 刷新）",
			path, len(committedLines), len(expectedLines))
	}

	return nil
}
