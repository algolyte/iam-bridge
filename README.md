# IAM-Bridge

IAM-Bridge is a flexible Identity and Access Management (IAM) wrapper service that provides a unified interface for integrating with various IAM providers such as Keycloak, Okta, Auth0, and AWS Cognito. The service allows switching between different IAM providers through configuration without requiring application redeployment.

## 🌟 Features

- **Multi-Provider Support**: Seamlessly integrate with different IAM providers
- **Hot Configuration Reload**: Change IAM providers without service restart
- **Standardized API**: Consistent API regardless of the underlying IAM provider
- **Security-First Design**: Implements security best practices
- **Observable**: Comprehensive logging and monitoring
- **Developer Friendly**: Hot reloading for rapid development
- **Docker Ready**: Containerized deployment support
- **Production Ready**: Includes health checks, graceful shutdown, and error handling

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for containerized deployment)
- Make (optional, for using Makefile commands)

## 🚀 Quick Start

1. **Clone the Repository**
   ```bash
   git clone https://github.com/zahidhasanpapon/iam-bridge
   cd iam-bridge
   ```

2. **Set Up Configuration**
   ```bash
   # Copy example environment file
   cp .env.example .env
   
   # Edit .env with your IAM provider credentials
   vim .env
   ```

3. **Run the Service**
   ```bash
   # Using Go directly
   go run cmd/api/main.go
   
   # Using Make
   make run
   
   # Using Docker Compose
   make docker-run
   ```

## 🛠️ Development

### Local Development
```bash
# Install development dependencies
make init

# Run with hot reload
make dev

# Run tests
make test

# Run linter
make lint
```

### Docker Development
```bash
# Build Docker image
make docker-build

# Run with Docker Compose
make docker-run

# Stop containers
make docker-stop
```

## 📝 Configuration

### Environment Variables

```env
# Application
APP_ENVIRONMENT=development
APP_PORT=8080

# Keycloak
KEYCLOAK_BASE_URL=http://localhost:8180
KEYCLOAK_REALM=master
KEYCLOAK_CLIENT_ID=your-client-id
KEYCLOAK_CLIENT_SECRET=your-client-secret

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json
```

### Configuration File (config/config.yaml)

```yaml
app:
  name: iam-bridge
  environment: development
  port: 8080
  debug: true

iam:
  provider: keycloak  # Can be: keycloak, okta, auth0, cognito
  keycloak:
    base_url: ${KEYCLOAK_BASE_URL}
    realm: ${KEYCLOAK_REALM}
    client_id: ${KEYCLOAK_CLIENT_ID}
    client_secret: ${KEYCLOAK_CLIENT_SECRET}

# ... additional configuration
```

## 🔌 API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Authenticate user
- `POST /api/v1/auth/logout` - Logout user
- `POST /api/v1/auth/refresh` - Refresh token
- `GET /api/v1/auth/validate` - Validate token

### User Management
- `GET /api/v1/users/:id` - Get user info
- `PUT /api/v1/users/:id` - Update user info
- `POST /api/v1/users/:id/roles` - Assign role
- `DELETE /api/v1/users/:id/roles/:role` - Remove role
- `GET /api/v1/users/:id/roles` - Get user roles

## 🔒 Security

- HTTPS/TLS support
- CORS configuration
- Rate limiting
- Request ID tracking
- Structured logging
- Panic recovery
- Error handling middleware

## 🏗️ Project Structure

```
iam-bridge/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── server/
│   ├── config/
│   ├── logger/
│   ├── middleware/
│   ├── providers/
│   └── routes/
├── pkg/
├── api/
├── config/
├── docs/
├── scripts/
├── test/
└── ... configuration files
```

## 🔧 Adding New IAM Providers

1. Create a new provider file in `internal/providers/`
2. Implement the `IAMProvider` interface
3. Add the provider to the factory in `iam_provider.go`
4. Update configuration structure in `config.go`

Example:
```go
type NewProvider struct {
    // Provider-specific fields
}

func NewNewProvider(cfg Config) (IAMProvider, error) {
    // Implementation
}

func (p *NewProvider) Login(ctx context.Context, username, password string) (string, error) {
    // Implementation
}

// Implement other interface methods
```

## 📊 Monitoring and Observability

- Structured JSON logging
- Request/Response logging
- Error tracking
- Health check endpoint
- Performance metrics

## 🚥 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run linter
make lint
```

## 📜 License

This project is licensed under the MIT License—see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## ✨ Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Viper](https://github.com/spf13/viper)

## 📬 Contact

Your Name - [@zahidhasan](https://www.linkedin.com/in/zahidhasanpapon/)

Project Link: [https://github.com/zahidhasanpapon/iam-bridge](https://github.com/zahidhasanpapon/iam-bridge)
