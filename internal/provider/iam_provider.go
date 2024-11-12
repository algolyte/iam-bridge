package provider

import (
	"context"
	"errors"
	"github.com/zahidhasanpapon/iam-bridge/internal/config"
	"github.com/zahidhasanpapon/iam-bridge/pkg/logger"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenInvalid       = errors.New("token invalid")
)

// TokenInfo represents the information extracted from a token
type TokenInfo struct {
	UserID    string                 `json:"user_id"`
	Username  string                 `json:"username"`
	Email     string                 `json:"email"`
	Roles     []string               `json:"roles"`
	Claims    map[string]interface{} `json:"claims"`
	ExpiresAt int64                  `json:"expires_at"`
}

// UserInfo represents the information of a user
type UserInfo struct {
	ID       string   `json:"id"`
	UserName string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

// IAMProvider defines the interface for all IAM providers must implement
type IAMProvider interface {
	Login(ctx context.Context, username, password string) (string, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (*TokenInfo, error)
	RefreshToken(ctx context.Context, token string) (string, error)

	GetUserInfo(ctx context.Context, userID string) (*UserInfo, error)
	UpdateUserInfo(ctx context.Context, userID string, userInfo *UserInfo) error

	AssignRole(ctx context.Context, userID, role string) error
	RemoveRole(ctx context.Context, userID, role string) error
	GetUserRoles(ctx context.Context, userID string) ([]string, error)

	HealthCheck(ctx context.Context) error
}

// NewIAMProvider creates a new IAM provider based on the given configuration
func NewIAMProvider(cfg *config.IAMConfig, log *logger.Logger) (IAMProvider, error) {
	switch cfg.CurrentProvider() {
	case "keycloak":
		return NewKeycloakProvider(cfg.Keycloak, log)
	default:
		return nil, errors.New("invalid IAM provider")
	}
}
