package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/zahidhasanpapon/iam-bridge/internal/config"
	"github.com/zahidhasanpapon/iam-bridge/internal/middleware"
	"github.com/zahidhasanpapon/iam-bridge/pkg/logger"
)

// ProviderSet is a provider set for wire
var ProviderSet = wire.NewSet(NewServer)

type Server struct {
	router     *gin.Engine
	config     *config.Config
	logger     logger.Logger
	httpServer *http.Server
}

func NewServer(cfg *config.Config, l logger.Logger) (*Server, error) {
	// Set Gin mode based on environment
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	s := &Server{
		router: router,
		config: cfg,
		logger: l,
	}

	// Setup middleware and routes
	s.setupMiddleware()
	s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	return s, nil
}

func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.router.Use(gin.Recovery())

	// Custom middleware
	s.router.Use(middleware.RequestID())
	s.router.Use(middleware.Logger(s.logger))
	s.router.Use(middleware.CORS(s.config))
	s.router.Use(middleware.ErrorHandler())
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}
}

func (s *Server) Start() error {
	s.logger.Infof("Starting server on port %s", s.config.ServerPort)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
