# Go Grafana Web Application

A modern Go web application with comprehensive monitoring using Grafana, Prometheus, and PostgreSQL. Built with clean architecture principles, dependency injection, and following Google Go standards.

## 🚀 Features

- **RESTful API**: Complete CRUD operations for user management
- **API Key Authentication**: Secure API key-based authentication for protected endpoints
- **API Key Management**: Full CRUD operations for managing API keys
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Dependency Injection**: Using Uber FX for clean dependency management
- **Monitoring**: Prometheus metrics collection and Grafana dashboards
- **Database**: PostgreSQL with GORM ORM
- **Containerization**: Docker and Kubernetes deployment ready
- **Documentation**: Swagger/OpenAPI documentation with API key support
- **Security**: CORS, input validation, secure headers, and API key validation
- **Logging**: Structured logging with Zap

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Layer    │    │  Business Layer │    │  Data Layer     │
│   (Handlers)    │◄──►│   (Services)    │◄──►│  (Repository)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Middleware    │    │   Domain        │    │   PostgreSQL    │
│   (Logging,     │    │   (Models)      │    │   Database      │
│    Metrics,     │    │                 │    │                 │
│    CORS)        │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📋 Prerequisites

- Go 1.24+
- Docker and Docker Compose
- PostgreSQL (for local development)
- Make (optional, for convenience)

## 🛠️ Installation & Setup

### Option 1: Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-grafana
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Verify services are running**
   ```bash
   docker-compose ps
   ```

### Option 2: Local Development

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL**
   ```bash
   # Using Docker
   docker run --name postgres -e POSTGRES_DB=go_grafana -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15-alpine
   ```
   ```sql
   INSERT INTO api_keys (
      name,
      key,
      description,
      active,
      expires_at,
      created_at,
      updated_at
   ) VALUES (
      'My API Key',
      'sk-123abctest',
      'API key for external service',
      true,
      '2024-12-31T23:59:59Z',
      NOW(),
      NOW()
   );
   ```

3. **Set environment variables**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=password
   export DB_NAME=go_grafana
   export DB_SSL_MODE=disable
   export SERVER_PORT=8080
   export LOG_LEVEL=info
   ```

4. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

## 🌐 API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

The API uses API key authentication for protected endpoints. Include your API key in the `X-API-Key` header:

```bash
curl -H "X-API-Key: sk-your-api-key-here" http://localhost:8080/api/v1/users
```

### User Management

| Method | Endpoint | Description | Authentication | Request Body |
|--------|----------|-------------|----------------|--------------|
| `POST` | `/users` | Create a new user | **Required** | `CreateUserRequest` |
| `GET` | `/users` | Get all users | Not required | - |
| `GET` | `/users/{id}` | Get user by ID | Not required | - |
| `PUT` | `/users/{id}` | Update user | **Required** | `UpdateUserRequest` |
| `DELETE` | `/users/{id}` | Delete user | **Required** | - |

### API Key Management

| Method | Endpoint | Description | Authentication | Request Body |
|--------|----------|-------------|----------------|--------------|
| `POST` | `/api-keys` | Create a new API key | **Required** | `CreateAPIKeyRequest` |
| `GET` | `/api-keys` | Get all API keys | **Required** | - |
| `GET` | `/api-keys/{id}` | Get API key by ID | **Required** | - |
| `PUT` | `/api-keys/{id}` | Update API key | **Required** | `UpdateAPIKeyRequest` |
| `DELETE` | `/api-keys/{id}` | Delete API key | **Required** | - |

### System Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/metrics` | Prometheus metrics |

### API Documentation

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **API Spec**: http://localhost:8080/swagger/doc.json

### Generating Swagger Documentation

To regenerate the Swagger documentation after making changes to the API:

```bash
# Run the generation script
./scripts/generate-swagger.sh

# Or manually with go run
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o docs
```

## 📊 Monitoring

### Grafana Dashboard
- **URL**: http://localhost:3000
- **Username**: `admin`
- **Password**: `admin`

### Prometheus
- **URL**: http://localhost:9090

### Available Metrics

#### HTTP Metrics
- `http_requests_total`: Total HTTP requests by method, endpoint, and status
- `http_request_duration_seconds`: Request duration histogram
- `http_requests_in_flight`: Current in-flight requests

#### Business Metrics
- `user_creation_total`: Total users created
- `user_deletion_total`: Total users deleted
- `user_update_total`: Total user updates
- `active_users_total`: Current active users count
- `user_age_distribution`: User age distribution histogram

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### API Testing Examples

#### Create API Key
```bash
curl -X POST http://localhost:8080/api/v1/api-keys \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk-your-existing-api-key" \
  -d '{
    "name": "My New API Key",
    "description": "API key for external service"
  }'
```

#### Get All API Keys
```bash
curl -H "X-API-Key: sk-your-api-key" http://localhost:8080/api/v1/api-keys
```

#### Create User (with API key)
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk-your-api-key" \
  -d '{
    "email": "john.doe@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "age": 30
  }'
```

#### Get All Users (no API key required)
```bash
curl http://localhost:8080/api/v1/users
```

#### Get User by ID (no API key required)
```bash
curl http://localhost:8080/api/v1/users/1
```

#### Update User (with API key)
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk-your-api-key" \
  -d '{
    "email": "john.doe.updated@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "age": 31,
    "active": true
  }'
```

#### Delete User (with API key)
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "X-API-Key: sk-your-api-key"
```

## 🔐 API Key Management

### Creating Your First API Key

Since all API key management endpoints require authentication, you'll need to create the first API key directly in the database or use a default key for initial setup.

#### Option 1: Database Insert
```sql
INSERT INTO api_keys (name, key, description, active, created_at, updated_at) 
VALUES ('Default API Key', 'sk-default-key-for-development', 'Default API key for development', true, NOW(), NOW());
```

#### Option 2: Environment Variable (for development)
You can modify the application to create a default API key on startup for development environments.

### API Key Security

- API keys are stored securely in the database
- Keys are hashed and validated on each request
- Expired or inactive keys are automatically rejected
- API keys can be set to expire at a specific date/time
- Keys are masked in API responses for security

### API Key Format

API keys follow the format: `sk-` followed by a 64-character hexadecimal string.

Example: `sk-1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef`

### API Key Management Examples

#### Create a New API Key
```bash
curl -X POST http://localhost:8080/api/v1/api-keys \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk-your-existing-api-key" \
  -d '{
    "name": "Production API Key",
    "description": "API key for production environment",
    "expires_at": "2024-12-31T23:59:59Z"
  }'
```

#### List All API Keys
```bash
curl -H "X-API-Key: sk-your-api-key" http://localhost:8080/api/v1/api-keys
```

#### Update an API Key
```bash
curl -X PUT http://localhost:8080/api/v1/api-keys/1 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk-your-api-key" \
  -d '{
    "name": "Updated API Key Name",
    "description": "Updated description",
    "active": true
  }'
```

#### Delete an API Key
```bash
curl -X DELETE http://localhost:8080/api/v1/api-keys/1 \
  -H "X-API-Key: sk-your-api-key"
```

## 🐳 Docker

### Build Image
```bash
docker build -f deployments/docker/Dockerfile -t go-grafana-app .
```

### Run Container
```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=password \
  -e DB_NAME=go_grafana \
  go-grafana-app
```

## ☸️ Kubernetes

### Deploy to Kubernetes
```bash
# Apply ConfigMap and Secret
kubectl apply -f deployments/k8s/configmap.yaml
kubectl create secret generic go-grafana-secret \
  --from-literal=db_user=postgres \
  --from-literal=db_password=password

# Deploy application
kubectl apply -f deployments/k8s/deployment.yaml
kubectl apply -f deployments/k8s/service.yaml
```

### Check Deployment
```bash
kubectl get pods
kubectl get services
kubectl logs -l app=go-grafana-app
```

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `password` | Database password |
| `DB_NAME` | `go_grafana` | Database name |
| `DB_SSL_MODE` | `disable` | Database SSL mode |
| `SERVER_PORT` | `8080` | Server port |
| `LOG_LEVEL` | `info` | Log level |

## 📁 Project Structure

```
go-grafana/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── domain/
│   │   ├── models/
│   │   │   ├── user.go            # Domain models
│   │   │   └── api_key.go         # API key models
│   │   └── repository/
│   │       ├── user_repository.go # Data access layer
│   │       └── api_key_repository.go # API key data access
│   ├── service/
│   │   ├── user_service.go        # Business logic
│   │   └── api_key_service.go     # API key business logic
│   ├── handler/
│   │   ├── user_handler.go        # HTTP handlers
│   │   └── api_key_handler.go     # API key HTTP handlers
│   └── middleware/
│       ├── logging.go             # Logging middleware
│       ├── metrics.go             # Metrics middleware
│       ├── cors.go                # CORS middleware
│       └── api_key_auth.go        # API key authentication
├── pkg/
│   ├── database/
│   │   └── postgres.go            # Database connection
│   └── metrics/
│       └── prometheus.go          # Custom metrics
├── deployments/
│   ├── docker/
│   │   └── Dockerfile             # Docker configuration
│   ├── k8s/
│   │   ├── deployment.yaml        # K8s deployment
│   │   ├── service.yaml           # K8s service
│   │   └── configmap.yaml         # K8s config
│   ├── prometheus/
│   │   └── prometheus.yml         # Prometheus config
│   └── grafana/
│       ├── dashboards/            # Grafana dashboards
│       └── datasources/           # Grafana datasources
├── docker-compose.yml             # Local development
├── go.mod                         # Go modules
├── go.sum                         # Go dependencies
└── README.md                      # This file
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/your-repo/go-grafana/issues) page
2. Create a new issue with detailed information
3. Contact the maintainers

## 🔄 Version History

- **v1.1.0**: Added API key authentication and management system
- **v1.0.0**: Initial release with basic CRUD operations and monitoring
- Future versions will include additional features and improvements

---

**Happy Coding! 🚀** 