package server

//
//import (
//	"github.com/zahidhasanpapon/iam-bridge/internal/config"
//	"github.com/zahidhasanpapon/iam-bridge/pkg/logger"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestServer_HealtCheck(t *testing.T) {
//	// Create test configuration
//	cfg := &config.Config{
//		AppName:    "test-app",
//		AppEnv:     "test",
//		ServerPort: "8080",
//		LogLevel:   "info",
//	}
//
//	// Create test logger
//	log, err := logger.NewLogger(cfg)
//	if err != nil {
//		t.Fatalf("failed to create logger: %v", err)
//	}
//
//	// Create test server instance
//	srv, err := NewServer(cfg, log)
//	if err != nil {
//		t.Fatalf("failed to create server: %v", err)
//	}
//	// Create test request
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/health", nil)
//
//	// Serve the request
//	srv.router.ServeHTTP(w, req)
//
//	// Assert response
//	if w.Code != http.StatusOK {
//		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
//	}
//
//	expectedBody := `{"status":"ok"}`
//	if w.Body.String() != expectedBody {
//		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
//	}
//}
