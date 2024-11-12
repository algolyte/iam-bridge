package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zahidhasanpapon/iam-bridge/internal/config"
	"github.com/zahidhasanpapon/iam-bridge/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type KeycloakProvider struct {
	config *config.KeycloakConfig
	logger *logger.Logger
	client *http.Client
}

// Login authenticates a user and returns an access token
func (k *KeycloakProvider) Login(ctx context.Context, username, password string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", k.config.ClientID)
	data.Set("client_secret", k.config.ClientSecret)
	data.Set("username", username)
	data.Set("password", password)

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		k.config.BaseURL, k.config.Realm)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL,
		strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := k.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to close response body: %w", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.AccessToken, nil
}

// ValidateToken validates the provided token and returns token information
func (k *KeycloakProvider) ValidateToken(ctx context.Context, token string) (*TokenInfo, error) {
	introspectionURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo",
		k.config.BaseURL, k.config.Realm)

	req, err := http.NewRequestWithContext(ctx, "GET", introspectionURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to close response body: %w", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrTokenInvalid
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var userInfo struct {
		Sub        string   `json:"sub"`
		Username   string   `json:"preferred_username"`
		Email      string   `json:"email"`
		RealmRoles []string `json:"realm_access"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &TokenInfo{
		UserID:   userInfo.Sub,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Roles:    userInfo.RealmRoles,
	}, nil
}

// Logout invalidates the provided token
func (k *KeycloakProvider) Logout(ctx context.Context, token string) error {
	logoutURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout",
		k.config.BaseURL, k.config.Realm)

	data := url.Values{}
	data.Set("client_id", k.config.ClientID)
	data.Set("client_secret", k.config.ClientSecret)
	data.Set("refresh_token", token)

	req, err := http.NewRequestWithContext(ctx, "POST", logoutURL,
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := k.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to close response body: %w", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// RefreshToken refreshes the provided token
func (k *KeycloakProvider) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", k.config.ClientID)
	data.Set("client_secret", k.config.ClientSecret)
	data.Set("refresh_token", refreshToken)

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		k.config.BaseURL, k.config.Realm)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL,
		strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := k.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to close response body: %w", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", ErrTokenExpired
		}
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.AccessToken, nil
}

// GetUserInfo retrieves user information
func (k *KeycloakProvider) GetUserInfo(ctx context.Context, userID string) (*UserInfo, error) {
	userURL := fmt.Sprintf("%s/admin/realms/%s/users/%s",
		k.config.BaseURL, k.config.Realm, userID)

	req, err := http.NewRequestWithContext(ctx, "GET", userURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to close response body: %w", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

func (k *KeycloakProvider) UpdateUserInfo(ctx context.Context, userID string, info *UserInfo) error {
	// Implementation details...
	return nil
}

func (k *KeycloakProvider) AssignRole(ctx context.Context, userID string, role string) error {
	// Implementation details...
	return nil
}

func (k *KeycloakProvider) RemoveRole(ctx context.Context, userID string, role string) error {
	// Implementation details...
	return nil
}

func (k *KeycloakProvider) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	// Implementation details...
	return nil, nil
}

func (k *KeycloakProvider) HealthCheck(ctx context.Context) error {
	healthURL := fmt.Sprintf("%s/health", k.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := k.client.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to close response body: %w", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// NewKeycloakProvider creates a new KeycloakProvider instance
func NewKeycloakProvider(cfg config.KeycloakConfig, log *logger.Logger) (IAMProvider, error) {
	if cfg.BaseURL == "" || cfg.Realm == "" || cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("missing required Keycloak configuration")
	}

	return &KeycloakProvider{
		config: &cfg,
		logger: log,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}
