package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type stubDatabase struct {
	err error
}

func (database stubDatabase) Ping(context.Context) error {
	return database.err
}

func testRouter(database DatabasePinger) stdhttp.Handler {
	return NewRouter(RouterOptions{
		Logger:   slog.New(slog.NewJSONHandler(io.Discard, nil)),
		Database: database,
		Version:  "test",
	}).Handler
}

func TestHealthReportsOKWhenDatabaseReachable(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, APIBasePath+"/health", nil)

	testRouter(stubDatabase{}).ServeHTTP(recorder, request)

	if recorder.Code != stdhttp.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, stdhttp.StatusOK)
	}

	var response healthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.State != "ok" || response.Database != "up" {
		t.Fatalf("response = %+v, want state ok and database up", response)
	}
}

func TestHealthReportsServiceUnavailableWhenDatabaseUnreachable(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, APIBasePath+"/health", nil)

	testRouter(stubDatabase{err: errors.New("connection refused")}).ServeHTTP(recorder, request)

	if recorder.Code != stdhttp.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", recorder.Code, stdhttp.StatusServiceUnavailable)
	}

	var response healthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.State != "degraded" || response.Database != "down" {
		t.Fatalf("response = %+v, want state degraded and database down", response)
	}
}

func TestOpenAPIDocumentKeepsRelativePathsAndServerPrefix(t *testing.T) {
	document := NewRouter(RouterOptions{
		Logger:  slog.New(slog.NewJSONHandler(io.Discard, nil)),
		Version: "test",
	}).API.OpenAPI()

	if _, found := document.Paths["/health"]; !found {
		t.Fatalf("openapi paths = %v, want /health", document.Paths)
	}

	if len(document.Servers) != 1 || document.Servers[0].URL != APIBasePath {
		t.Fatalf("openapi servers = %+v, want single server %s", document.Servers, APIBasePath)
	}
}

func TestRequestLogUsesForwardedClientIPWhenProxyIsTrusted(t *testing.T) {
	var logOutput bytes.Buffer
	router := NewRouter(RouterOptions{
		Logger:            slog.New(slog.NewJSONHandler(&logOutput, nil)),
		Version:           "test",
		TrustedProxyCount: 1,
	})

	request := httptest.NewRequest(stdhttp.MethodGet, APIBasePath+"/health", nil)
	request.RemoteAddr = "10.0.0.1:5432"
	request.Header.Set("X-Forwarded-For", "203.0.113.9")
	router.Handler.ServeHTTP(httptest.NewRecorder(), request)

	if !strings.Contains(logOutput.String(), `"client_ip":"203.0.113.9"`) {
		t.Fatalf("log output = %s, want client_ip 203.0.113.9", logOutput.String())
	}
}

func TestRequestLogIgnoresForwardedClientIPWithoutTrustedProxy(t *testing.T) {
	var logOutput bytes.Buffer
	router := NewRouter(RouterOptions{
		Logger:  slog.New(slog.NewJSONHandler(&logOutput, nil)),
		Version: "test",
	})

	request := httptest.NewRequest(stdhttp.MethodGet, APIBasePath+"/health", nil)
	request.RemoteAddr = "10.0.0.1:5432"
	request.Header.Set("X-Forwarded-For", "203.0.113.9")
	router.Handler.ServeHTTP(httptest.NewRecorder(), request)

	if !strings.Contains(logOutput.String(), `"client_ip":"10.0.0.1"`) {
		t.Fatalf("log output = %s, want client_ip 10.0.0.1", logOutput.String())
	}
}

func TestRequestIDIgnoresInboundHeaderAndIsEchoedBack(t *testing.T) {
	var logOutput bytes.Buffer
	router := NewRouter(RouterOptions{
		Logger:  slog.New(slog.NewJSONHandler(&logOutput, nil)),
		Version: "test",
	})

	request := httptest.NewRequest(stdhttp.MethodGet, APIBasePath+"/health", nil)
	request.Header.Set(requestIDHeaderName, "forged-by-client")
	recorder := httptest.NewRecorder()
	router.Handler.ServeHTTP(recorder, request)

	requestID := recorder.Header().Get(requestIDHeaderName)
	if requestID == "" {
		t.Fatal("response header is empty, want a generated request id")
	}
	if requestID == "forged-by-client" {
		t.Fatal("inbound request id was reused, want a generated one")
	}
	if !strings.Contains(logOutput.String(), `"request_id":"`+requestID+`"`) {
		t.Fatalf("log output = %s, want request_id %s", logOutput.String(), requestID)
	}
}

func TestRecoveringTurnsPanicIntoProblemResponseAndLogsIt(t *testing.T) {
	var logOutput bytes.Buffer
	panicHandler := recovering(slog.New(slog.NewJSONHandler(&logOutput, nil)))(
		stdhttp.HandlerFunc(func(stdhttp.ResponseWriter, *stdhttp.Request) {
			panic("boom")
		}),
	)

	recorder := httptest.NewRecorder()
	panicHandler.ServeHTTP(recorder, httptest.NewRequest(stdhttp.MethodGet, "/any", nil))

	if recorder.Code != stdhttp.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, stdhttp.StatusInternalServerError)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/problem+json" {
		t.Fatalf("content type = %q, want application/problem+json", contentType)
	}

	var document internalErrorDocument
	if err := json.Unmarshal(recorder.Body.Bytes(), &document); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if document.Status != stdhttp.StatusInternalServerError {
		t.Fatalf("document = %+v, want status 500", document)
	}

	logged := logOutput.String()
	if !strings.Contains(logged, `"panic":"boom"`) || !strings.Contains(logged, `"stack":`) {
		t.Fatalf("log output = %s, want panic and stack fields", logged)
	}
}
