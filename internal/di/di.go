package di

import (
	"github.com/google/wire"
	"github.com/zahidhasanpapon/iam-bridge/internal/config"
	"github.com/zahidhasanpapon/iam-bridge/internal/server"
	"github.com/zahidhasanpapon/iam-bridge/pkg/logger"
)

// ProviderSet is a provider set for wire
var ProviderSet = wire.NewSet(
	config.ProviderSet,
	logger.ProviderSet,
	server.ProviderSet,
)

// InitializeApp initializes the application
func InitializeApp() (*server.Server, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
