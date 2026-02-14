# Golang DDD Template

Production-ready REST API template built with **Clean Architecture** and **Domain-Driven Design (DDD)** principles using Go.

## 🚀 Features

- ✅ **Clean Architecture** with DDD structure
- ✅ **Production-ready** code (no placeholders)
- ✅ **OWASP Top 10** security compliance
- ✅ **PostgreSQL** with pgx driver
- ✅ **Redis** caching
- ✅ **RabbitMQ** message queue
- ✅ **JWT** authentication
- ✅ **Comprehensive validation** (go-playground/validator)
- ✅ **Structured logging** (Zap)
- ✅ **Swagger documentation**
- ✅ **Docker & Docker Compose**
- ✅ **Database migrations**
- ✅ **Prometheus & Grafana** monitoring
- ✅ **Rate limiting** & **CORS**
- ✅ **Graceful shutdown**
- ✅ **Example User/Auth domain**

## 📁 Project Structure

```
.
├── cmd/
│   └── api/                    # Application entry points
│       └── main.go
├── internal/
│   ├── domain/                 # Domain layer (Business logic)
│   │   ├── user/
│   │   │   ├── entity/        # Domain entities
│   │   │   ├── repository/    # Repository interfaces & implementations
│   │   │   ├── usecase/       # Business use cases
│   │   │   ├── dto/           # Data Transfer Objects
│   │   │   └── delivery/      # HTTP handlers
│   │   └── auth/              # Auth domain (similar structure)
│   ├── infrastructure/         # Infrastructure layer
│   │   ├── database/          # Database connections
│   │   ├── cache/             # Redis cache
│   │   ├── messaging/         # RabbitMQ
│   │   └── config/            # Configuration management
│   ├── delivery/              # Delivery layer
│   │   └── http/
│   │       ├── middleware/    # HTTP middlewares
│   │       ├── handler/       # HTTP handlers
│   │       └── router/        # Route definitions
│   └── shared/                # Shared utilities
│       ├── errors/            # Custom errors
│       ├── constants/         # Constants
│       └── utils/             # Utility functions
├── pkg/                       # Public packages
│   ├── logger/               # Logging utility
│   ├── validator/            # Validation utility
│   ├── response/             # HTTP response utility
│   ├── jwt/                  # JWT utility
│   └── crypto/               # Cryptography utility
├── migrations/               # Database migrations
├── docs/                     # Swagger documentation
├── tests/                    # Tests
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── deployments/              # Deployment configs
│   └── docker/
├── scripts/                  # Utility scripts
├── .env.example             # Environment variables template
├── docker-compose.yml       # Docker Compose configuration
├── Dockerfile               # Docker image definition
├── Makefile                 # Development commands
└── README.md

```

## 🛠 Tech Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | 1.21 |
| Framework | Gin | 1.10.0 |
| Database | PostgreSQL | 16 |
| Database Driver | pgx | 5.5.0 |
| Cache | Redis | 7 |
| Message Queue | RabbitMQ | 3.12 |
| Migration | golang-migrate | 4.17.0 |
| Authentication | JWT | 5.2.0 |
| Validation | validator | 10.16.0 |
| Configuration | Viper | 1.18.2 |
| Logging | Zap | 1.26.0 |
| Monitoring | Prometheus + Grafana | Latest |
| Documentation | Swagger | 1.16.2 |

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional, for easier commands)
- golang-migrate CLI (for migrations)

### 1. Clone the repository

```bash
git clone https://github.com/TubagusAldiMY/go-template.git
cd golang-ddd-template
```

### 2. Copy environment file

```bash
cp .env.example .env
```

### 3. Start with Docker Compose (Recommended)

```bash
# Start all services (PostgreSQL, Redis, RabbitMQ, App)
make docker-up

# Run migrations
make migrate-up

# View logs
make docker-logs
```

The API will be available at `http://localhost:8080`

### 4. Or run locally

```bash
# Install dependencies
make deps

# Start infrastructure services only
docker-compose up -d postgres redis rabbitmq

# Run migrations
make migrate-up

# Generate Swagger docs
make swagger

# Run application
make run
```

## 📚 API Documentation

Once the application is running, access Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## 🔐 Authentication

The API uses JWT Bearer tokens for authentication.

### Register a new user

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

### Use the token

```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <your-access-token>"
```

## 🗄 Database Migrations

```bash
# Create a new migration
make migrate-create name=create_products_table

# Run migrations
make migrate-up

# Rollback last migration
make migrate-down

# Force migration version
make migrate-force version=1
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run unit tests only
make test-unit
```

## 🛠 Development Commands

```bash
# View all available commands
make help

# Format code
make fmt

# Run linter
make lint

# Generate Swagger docs
make swagger

# Install development tools
make install-tools

# Complete setup (first time)
make setup
```

## 📊 Monitoring

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)

## 🏗 Architecture

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`)
   - Entities: Core business objects
   - Repository Interfaces: Data access contracts
   - Use Cases: Business logic
   - DTOs: Data transfer objects

2. **Infrastructure Layer** (`internal/infrastructure/`)
   - Database implementations
   - Cache implementations
   - Message queue implementations
   - External service integrations

3. **Delivery Layer** (`internal/delivery/`)
   - HTTP handlers
   - Middleware
   - Request/Response mapping

4. **Shared Layer** (`internal/shared/`)
   - Common utilities
   - Error definitions
   - Constants

### Dependency Rule

Dependencies point inward:
```
Delivery → Use Cases → Entities
Infrastructure → Use Cases → Entities
```

## 🔒 Security Features

✅ **Password Security**
- Bcrypt hashing (configurable cost)
- Password strength validation

✅ **JWT Security**
- HS256 signing
- Token expiration
- Refresh token support

✅ **HTTP Security**
- CORS configuration
- Rate limiting
- Request ID tracking
- Secure headers

✅ **Input Validation**
- Request validation
- SQL injection prevention (pgx parameterized queries)
- XSS protection

✅ **OWASP Top 10 Compliance**
- A01: Broken Access Control → Role-based access control
- A02: Cryptographic Failures → Bcrypt + JWT
- A03: Injection → Parameterized queries
- A04: Insecure Design → Clean Architecture
- A05: Security Misconfiguration → Environment-based config
- A06: Vulnerable Components → Regular dependency updates
- A07: Authentication Failures → JWT + password policies
- A08: Software/Data Integrity → Code signing, migrations
- A09: Logging Failures → Structured logging with Zap
- A10: SSRF → Input validation

## 📝 Adding a New Domain

1. **Create domain structure**:
```bash
mkdir -p internal/domain/product/{entity,repository,usecase,dto,delivery/http}
```

2. **Define entity** (`internal/domain/product/entity/product.go`)
3. **Define repository interface** (`internal/domain/product/repository/product_repository.go`)
4. **Implement repository** (`internal/domain/product/repository/postgres_product_repository.go`)
5. **Define DTOs** (`internal/domain/product/dto/product_dto.go`)
6. **Implement use cases** (`internal/domain/product/usecase/product_usecase.go`)
7. **Create handlers** (`internal/domain/product/delivery/http/product_handler.go`)
8. **Register routes** in `internal/delivery/http/router/router.go`
9. **Create migration** for the new table

## 🔄 CI/CD (To be added)

This template is ready for CI/CD integration. Add your preferred pipeline:
- GitHub Actions
- GitLab CI
- Jenkins
- CircleCI

## 📄 License

This project is licensed under the MIT License.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📧 Support

For issues and questions, please open an issue in the repository.

---

**Happy Coding! 🚀**
