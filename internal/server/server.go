package server

import (
	"context"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zahidhasanpapon/iam-bridge/internal/config"
	"github.com/zahidhasanpapon/iam-bridge/internal/middleware"
	"github.com/zahidhasanpapon/iam-bridge/internal/provider"
	"github.com/zahidhasanpapon/iam-bridge/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server and its dependencies
type Server struct {
	config      *config.Config
	logger      logger.Logger
	router      *gin.Engine
	iamProvider provider.IAMProvider
	httpServer  *http.Server
}

// NewServer creates a new server instance
func NewServer() (*Server, error) {
	// Load configuration
	cfg, err := config.LoadConfig("config")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	log, _ := logger.NewLogger(&cfg.Logging)

	// Initialize IAM provider
	iamProvider, err := provider.NewIAMProvider(&cfg.IAM, &log)
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM provider: %w", err)
	}

	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router with default middleware
	router := gin.New()

	// Create server instance
	server := &Server{
		config:      cfg,
		logger:      log,
		router:      router,
		iamProvider: iamProvider,
	}

	// Initialize server
	server.setupMiddleware()
	server.setupRoutes()

	return server, nil
}

// setupMiddleware configures all middleware for the server
func (s *Server) setupMiddleware() {
	// Swagger endpoint
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Add basic middleware
	s.router.Use(
		middleware.RequestIDMiddleware(),
		middleware.LoggerMiddleware(s.logger),
		middleware.RecoveryMiddleware(s.logger),
		middleware.CORSMiddleware(&s.config.Security.CORS),
		middleware.ErrorHandlerMiddleware(),
	)

	// Add rate limiting if enabled
	if s.config.Security.RateLimit.Enabled {
		// TODO: Implement rate limiting middleware
	}
}

// setupRoutes configures all routes for the server
func (s *Server) setupRoutes() {
	// Health check
	// @Summary Health check
	// @Description Checks the health of the service
	// @Tags Health
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 503 {object} map[string]interface{}
	// @Router /health [get]
	s.router.GET("/health", s.handleHealthCheck)

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			// @Summary Login
			// @Description Authenticates a user and provides a token
			// @Tags Authentication
			// @Accept json
			// @Produce json
			// @Param credentials body struct{Username string; Password string} true "Login credentials"
			// @Success 200 {object} map[string]string
			// @Failure 400 {object} map[string]interface{}
			// @Router /api/v1/auth/login [post]
			auth.POST("/login", s.handleLogin)

			// @Summary Logout
			// @Description Logs out a user by invalidating their token
			// @Tags Authentication
			// @Security BearerAuth
			// @Success 204
			// @Failure 400 {object} map[string]interface{}
			// @Router /api/v1/auth/logout [post]
			auth.POST("/logout", s.handleLogout)

			// @Summary Refresh Token
			// @Description Refreshes an access token using a refresh token
			// @Tags Authentication
			// @Accept json
			// @Produce json
			// @Param refreshToken body struct{RefreshToken string} true "Refresh token"
			// @Success 200 {object} map[string]string
			// @Failure 400 {object} map[string]interface{}
			// @Router /api/v1/auth/refresh [post]
			auth.POST("/refresh", s.handleRefreshToken)

			// @Summary Validate Token
			// @Description Validates the provided token
			// @Tags Authentication
			// @Security BearerAuth
			// @Success 200 {object} map[string]interface{}
			// @Failure 400 {object} map[string]interface{}
			// @Router /api/v1/auth/validate [get]
			auth.GET("/validate", s.handleValidateToken)
		}

		// User management routes
		users := api.Group("/users")
		{
			// @Summary Get User Info
			// @Description Retrieves information about a specific user
			// @Tags Users
			// @Param id path string true "User ID"
			// @Produce json
			// @Success 200 {object} map[string]interface{}
			// @Failure 400 {object} map[string]interface{}
			// @Router /api/v1/users/{id} [get]
			users.GET("/:id", s.handleGetUserInfo)

			// @Summary Update User Info
			// @Description Updates the information of a specific user
			// @Tags Users
			// @Param id path string true "User ID"
			// @Param userInfo body struct{...} true "User information payload"
			// @Success 204
			// @Failure 400 {object} map[string]interface{}
			// @Router /api/v1/users/{id} [put]
			users.PUT("/:id", s.handleUpdateUserInfo)

			// @Summary Assign Role
			// @Description Assigns a role to a specific user
			// @Tags Users
			// @Param id path string true "User ID"
			// @Param role body struct{Role string} true "Role payload"
			// @Success 204
			// @Failure 400 {object} map[string]interface{}
			// @Router /api/v1/users/{id}/roles [post]
			users.POST("/:id/roles", s.handleAssignRole)

			// @Summary Remove Role
			// @Description Removes a role from a specific user
			// @Tags Users
			// @Param id path string true "User ID"
			// @Param role path string true "Role"
			// @Success 204
			// @Failure 400 {object} map[string]interface{}
			// @Router /api/v1/users/{id}/roles/{role} [delete]
			users.DELETE("/:id/roles/:role", s.handleRemoveRole)

			// @Summary Get User Roles
			// @Description Retrieves the roles of a specific user
			// @Tags Users
			// @Param id path string true "User ID"
			// @Produce json
			// @Success 200 {object} map[string]interface{}
			// @Failure 400 {object} map[string]interface{}
			// @Router /api/v1/users/{id}/roles [get]
			users.GET("/:id/roles", s.handleGetUserRoles)
		}
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.App.Port),
		Handler: s.router,
	}

	// Create a channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		s.logger.Info("Starting server", "port", s.config.App.Port)
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	// Create a channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or an error from the server
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case <-shutdown:
		s.logger.Info("Starting shutdown")

		// Create a deadline for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Shut down the server
		if err := s.httpServer.Shutdown(ctx); err != nil {
			// If shutdown times out, forcefully close
			s.httpServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

// Handler functions
func (s *Server) handleHealthCheck(c *gin.Context) {
	// Check IAM provider health
	if err := s.iamProvider.HealthCheck(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "IAM provider health check failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
	})
}

func (s *Server) handleLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("invalid request: %w", err))
		return
	}

	token, err := s.iamProvider.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (s *Server) handleLogout(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.Error(provider.ErrTokenInvalid)
		return
	}

	if err := s.iamProvider.Logout(c.Request.Context(), token); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) handleRefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("invalid request: %w", err))
		return
	}

	token, err := s.iamProvider.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (s *Server) handleValidateToken(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.Error(provider.ErrTokenInvalid)
		return
	}

	tokenInfo, err := s.iamProvider.ValidateToken(c.Request.Context(), token)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tokenInfo)
}

func (s *Server) handleGetUserInfo(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.Error(fmt.Errorf("user ID is required"))
		return
	}

	userInfo, err := s.iamProvider.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

func (s *Server) handleUpdateUserInfo(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.Error(fmt.Errorf("user ID is required"))
		return
	}

	var userInfo provider.UserInfo
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.Error(fmt.Errorf("invalid request: %w", err))
		return
	}

	if err := s.iamProvider.UpdateUserInfo(c.Request.Context(), userID, &userInfo); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) handleAssignRole(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.Error(fmt.Errorf("user ID is required"))
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("invalid request: %w", err))
		return
	}

	if err := s.iamProvider.AssignRole(c.Request.Context(), userID, req.Role); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) handleRemoveRole(c *gin.Context) {
	userID := c.Param("id")
	role := c.Param("role")
	if userID == "" || role == "" {
		c.Error(fmt.Errorf("user ID and role are required"))
		return
	}

	if err := s.iamProvider.RemoveRole(c.Request.Context(), userID, role); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) handleGetUserRoles(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.Error(fmt.Errorf("user ID is required"))
		return
	}

	roles, err := s.iamProvider.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
	})
}

// Helper functions
func extractToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		return ""
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		return token[7:]
	}

	return token
}
