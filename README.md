# Go AI Security

A modern Go application built with Domain-Driven Design (DDD) and Clean Architecture principles, designed for AI security applications.

## 🏗️ Architecture

This project follows **Domain-Driven Design (DDD)** and **Clean Architecture** patterns, inspired by NestJS module philosophy:

- **Domain Layer** → Pure business rules and entities
- **UseCase Layer** → Business logic orchestration  
- **Repository Layer** → Data persistence (MongoDB, Redis, etc.)
- **Delivery Layer** → API handlers (Gin), message consumers
- **DTOs** → Input/output models for external communication

## 📁 Project Structure

```
/cmd/                    # Application entry point
/config/                 # Configuration management
/internal/               # Internal modules
  /users/               # Users module
    /domain/            # Domain entities, errors, events
    /usecase/           # Application services
    /repository/        # Repository interfaces & implementations
    /delivery/http/     # HTTP handlers and routes
    /dto/               # Request/Response DTOs
  /auth/                # Authentication module
    /domain/            # Domain entities, errors, events
    /usecase/           # Application services
    /repository/        # Repository interfaces & implementations
    /delivery/http/     # HTTP handlers and routes
    /dto/               # Request/Response DTOs
/pkg/                   # Shared utilities
  /logger/              # Centralized logging
  /middleware/          # HTTP middleware
  /shared/              # Shared types and constants
  /utils/               # Utility functions
```

## 🚀 Getting Started

### Prerequisites

- Go 1.25.3 or higher
- MongoDB (for data persistence)
- Redis (for caching and sessions)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-ai-security
```

2. Install dependencies:
```bash
go mod download
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

### Running the Application

```bash
# Run the application
go run cmd/main.go

# Or build and run
go build -o bin/go-ai-security cmd/main.go
./bin/go-ai-security
```

## 🛠️ Development

### Code Quality

The project includes several tools for maintaining code quality:

- **golangci-lint**: Comprehensive Go linting
- **Prettier**: Code formatting
- **EditorConfig**: Editor configuration

### Available Commands

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Build application
go build -o bin/go-ai-security cmd/main.go
```

### Code Style Guidelines

- Use tabs for indentation in Go files
- Maximum line length: 120 characters
- Follow Go naming conventions
- Write comprehensive tests
- Use meaningful variable and function names
- Keep functions small and focused
- Use interfaces for dependency injection

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run tests with verbose output
go test -v ./...
```

## 📦 Modules

### Users Module
Handles user management functionality including:
- User registration and authentication
- User profile management
- User data persistence

### Auth Module
Manages authentication and authorization:
- JWT token generation and validation
- Session management
- Access control

## 🔧 Configuration

Configuration is managed through environment variables and config files in the `/config` directory.

### Environment Variables

- `GIN_MODE`: Gin mode (debug/release)
- `PORT`: Application port (default: 8080)
- `MONGODB_URI`: MongoDB connection string
- `REDIS_URL`: Redis connection string
- `JWT_SECRET`: JWT signing secret

## 📚 API Documentation

API endpoints are documented using Swagger/OpenAPI specifications. Once the application is running, you can access the API documentation at:

```
http://localhost:8080/swagger/index.html
```

## 🏛️ Design Principles

### Domain-Driven Design (DDD)
- **Bounded Contexts**: Each module represents a distinct business domain
- **Entities**: Core business objects with identity
- **Value Objects**: Immutable objects without identity
- **Domain Events**: Business events that occur in the domain

### Clean Architecture
- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Interface Segregation**: Small, focused interfaces
- **Single Responsibility**: Each layer has a single responsibility
- **Dependency Injection**: Dependencies are injected, not created

### SOLID Principles
- **S**ingle Responsibility Principle
- **O**pen/Closed Principle
- **L**iskov Substitution Principle
- **I**nterface Segregation Principle
- **D**ependency Inversion Principle

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow the existing code style
- Write tests for new functionality
- Update documentation as needed
- Ensure all tests pass
- Run linter before committing

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support, please open an issue in the repository or contact the development team.

## 🔗 Related Projects

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [Go Redis](https://github.com/redis/go-redis)

## 📈 Roadmap

- [ ] Add more authentication providers
- [ ] Implement rate limiting
- [ ] Add comprehensive logging
- [ ] Implement caching strategies
- [ ] Add monitoring and metrics
- [ ] Implement message queuing
- [ ] Add comprehensive API documentation

---

**Built with ❤️ using Go and Clean Architecture principles**
