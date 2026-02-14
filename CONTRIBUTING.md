# Contributing to Golang DDD Template

Thank you for your interest in contributing! This document provides guidelines for contributing to this project.

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Respect different viewpoints and experiences

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in Issues
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, etc.)
   - Error messages or logs

### Suggesting Features

1. Check if the feature has been suggested
2. Create an issue describing:
   - The problem it solves
   - Proposed solution
   - Alternative solutions considered
   - Impact on existing functionality

### Pull Requests

1. **Fork and Clone**
   ```bash
   git clone https://github.com/TubagusAldiMY/go-template.git
   cd golang-ddd-template
   ```

2. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Changes**
   - Follow the coding standards (see below)
   - Write tests for new functionality
   - Update documentation if needed
   - Run tests and linter

4. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add amazing feature"
   ```

   Use conventional commit messages:
   - `feat:` - New feature
   - `fix:` - Bug fix
   - `docs:` - Documentation changes
   - `style:` - Code style changes (formatting)
   - `refactor:` - Code refactoring
   - `test:` - Adding or updating tests
   - `chore:` - Maintenance tasks

5. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   ```

## Coding Standards

### Go Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Keep functions small and focused
- Write clear, descriptive names

### Clean Architecture Principles

- **Dependency Rule**: Dependencies point inward
  - Domain layer has no external dependencies
  - Use cases depend only on domain
  - Infrastructure implements interfaces from domain

- **Layer Separation**:
  ```
  Domain (entities, interfaces) ← Use Cases ← Infrastructure
                                            ← Delivery
  ```

### DDD Principles

- **Entities**: Core business objects with identity
- **Value Objects**: Immutable objects without identity
- **Repositories**: Data access abstraction
- **Use Cases**: Business logic orchestration
- **DTOs**: Data transfer between layers

### Testing

- Write unit tests for use cases
- Write integration tests for repositories
- Aim for >80% code coverage
- Use table-driven tests when appropriate
- Mock external dependencies

Example:
```go
func TestUserUsecase_Register(t *testing.T) {
    tests := []struct {
        name    string
        input   *dto.RegisterRequest
        want    *dto.UserResponse
        wantErr error
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Documentation

- Document all exported functions
- Use godoc format
- Include examples for complex functions
- Update README for significant changes
- Add Swagger annotations for API endpoints

### Security

- Never commit secrets or credentials
- Validate all user inputs
- Use parameterized queries
- Hash passwords with bcrypt
- Implement rate limiting
- Follow OWASP guidelines

## Development Workflow

1. **Setup Development Environment**
   ```bash
   make setup
   ```

2. **Run Tests**
   ```bash
   make test
   ```

3. **Format Code**
   ```bash
   make fmt
   ```

4. **Run Linter**
   ```bash
   make lint
   ```

5. **Generate Swagger Docs**
   ```bash
   make swagger
   ```

## Adding a New Domain

When adding a new domain (e.g., Product):

1. Create directory structure:
   ```bash
   mkdir -p internal/domain/product/{entity,repository,usecase,dto,delivery/http}
   ```

2. Create entity:
   ```go
   // internal/domain/product/entity/product.go
   package entity
   
   type Product struct {
       ID    string
       Name  string
       // ...
   }
   ```

3. Define repository interface:
   ```go
   // internal/domain/product/repository/product_repository.go
   package repository
   
   type ProductRepository interface {
       Create(ctx context.Context, product *entity.Product) error
       // ...
   }
   ```

4. Implement repository:
   ```go
   // internal/domain/product/repository/postgres_product_repository.go
   ```

5. Create DTOs:
   ```go
   // internal/domain/product/dto/product_dto.go
   ```

6. Implement use cases:
   ```go
   // internal/domain/product/usecase/product_usecase.go
   ```

7. Create HTTP handlers:
   ```go
   // internal/domain/product/delivery/http/product_handler.go
   ```

8. Register routes in router

9. Create database migration

10. Write tests

## Review Process

1. All PRs require at least one review
2. CI/CD checks must pass
3. Code coverage should not decrease
4. Documentation must be updated
5. Tests must be included

## Questions?

Feel free to open an issue with your question or reach out to maintainers.

Thank you for contributing! 🎉
