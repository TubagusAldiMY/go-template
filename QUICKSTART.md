# Quick Start Guide

Get up and running with Golang DDD Template in 5 minutes!

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional but recommended)

## 🚀 5-Minute Setup

### 1. Clone & Setup

```bash
git clone https://github.com/TubagusAldiMY/go-template.git
cd golang-ddd-template
cp .env.example .env
```

### 2. Start Services

```bash
make docker-up
```

This starts:
- PostgreSQL (port 5432)
- Redis (port 6379)
- RabbitMQ (port 5672, Management: 15672)
- Application (port 8080)
- Prometheus (port 9090)
- Grafana (port 3000)

### 3. Run Migrations

```bash
make migrate-up
```

### 4. Seed Database (Optional)

```bash
make seed
```

This creates:
- Admin user: `admin@example.com` / `Admin123!`
- Test user: `user@example.com` / `User123!`

### 5. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Register new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "username": "johndoe",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

## 📚 Next Steps

### View API Documentation

Open Swagger UI:
```
http://localhost:8080/swagger/index.html
```

### Access Monitoring

- **RabbitMQ Management**: http://localhost:15672 (guest/guest)
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

### Development Commands

```bash
# View all available commands
make help

# Run application locally
make run

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Generate Swagger docs
make swagger

# View Docker logs
make docker-logs

# Stop all services
make docker-down
```

## 📝 Using the Template

### Customize for Your Project

1. **Update module name** in `go.mod`:
   ```
   module github.com/yourorg/yourproject
   ```

2. **Find and replace** all occurrences:
   ```bash
   find . -type f -name "*.go" -exec sed -i 's/github.com\/yourusername\/golang-ddd-template/github.com\/yourorg\/yourproject/g' {} +
   ```

3. **Update configuration** in `.env`:
   - Change `APP_NAME`
   - Set secure `JWT_SECRET`
   - Configure database credentials

4. **Update Swagger info** in `cmd/api/main.go`

### Add a New Feature

1. **Create domain structure**:
   ```bash
   mkdir -p internal/domain/product/{entity,repository,usecase,dto,delivery/http}
   ```

2. **Follow the User domain** as example

3. **Create migration**:
   ```bash
   make migrate-create name=create_products_table
   ```

4. **Update router** to add new routes

5. **Run tests**:
   ```bash
   make test
   ```

## 🔒 Security

Before production:

1. Change default passwords
2. Use strong JWT secret (32+ random characters)
3. Enable HTTPS/TLS
4. Configure CORS properly
5. Review and adjust rate limits
6. Disable debug mode (`APP_DEBUG=false`)
7. Use environment-specific `.env` files

## 🐛 Troubleshooting

### Port already in use

```bash
# Stop all services
make docker-down

# Check ports
lsof -i :8080
lsof -i :5432

# Kill process using port
kill -9 <PID>
```

### Database connection failed

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# View logs
docker logs ddd-postgres

# Restart services
make docker-down
make docker-up
```

### Migration failed

```bash
# Check migration status
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/ddd_template?sslmode=disable" version

# Force specific version
make migrate-force version=1

# Then run migrations again
make migrate-up
```

## 📖 Learn More

- [Full README](README.md)
- [Contributing Guide](CONTRIBUTING.md)
- [Security Policy](SECURITY.md)
- [Changelog](CHANGELOG.md)

## 💡 Tips

- Use `make help` to see all available commands
- Check Swagger docs for complete API reference
- Review example tests in `tests/unit/`
- Follow User domain pattern for new features
- Keep dependencies updated regularly

## 🎉 You're Ready!

Start building your awesome API! 🚀

For questions or issues, please check the [README](README.md) or open an issue.
