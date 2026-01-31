# Go Clean Architecture Boilerplate

A production-ready Go project template implementing Clean Architecture with essential libraries and best practices.

## 🚀 Features

- **Clean Architecture** - Separation of concerns with Entity, Repository, UseCase, and Handler layers
- **Gin Framework** - Fast HTTP web framework
- **GORM** - Full-featured ORM for PostgreSQL
- **JWT Authentication** - Secure token-based authentication
- **Redis** - Caching and session management
- **Swagger/OpenAPI** - API documentation with swaggo
- **Database Migrations** - Version control for database schema
- **Docker Support** - Development and production containers
- **Hot Reload** - Air for fast development
- **Structured Logging** - Logrus for logging
- **Configuration** - Viper for environment management
- **Validation** - Go Playground Validator
- **Email Service** - Gomail for SMTP

## 📁 Project Structure

```
go-clean-architecture/
├── cmd/
│   ├── api/                    # Main application entry point
│   │   └── main.go
│   └── migrate/                # Database migration CLI
│       └── main.go
├── config/
│   └── config.go               # Configuration management
├── database/
│   ├── migrations/             # SQL migration files
│   └── seeder/                 # Database seeders
├── docs/                       # Swagger documentation
├── internal/
│   ├── dto/                    # Data Transfer Objects
│   ├── entity/                 # Domain entities
│   ├── handler/                # HTTP handlers (controllers)
│   ├── middleware/             # HTTP middleware
│   ├── repository/             # Data access layer
│   ├── router/                 # Route definitions
│   └── usecase/                # Business logic layer
├── pkg/
│   ├── database/               # Database connections
│   ├── logger/                 # Logging utilities
│   ├── mail/                   # Email service
│   ├── response/               # HTTP response helpers
│   ├── utils/                  # Utility functions
│   └── validator/              # Validation helpers
├── .air.toml                   # Air hot reload config
├── .env.example                # Environment template
├── .gitignore
├── docker-compose.yml          # Production Docker setup
├── docker-compose.dev.yml      # Development Docker setup
├── Dockerfile                  # Production Dockerfile
├── Dockerfile.dev              # Development Dockerfile
├── go.mod
├── Makefile                    # Development commands
└── README.md
```

## 🛠️ Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional, for Makefile commands)
- PostgreSQL (if running locally)
- Redis (if running locally)

## 🚀 Quick Start

### 1. Clone and Setup

```bash
# Clone the repository
git clone https://github.com/your-username/go-clean-architecture.git
cd go-clean-architecture

# Copy environment file
cp .env.example .env

# Edit .env with your configuration
```

### 2. Install Development Tools

```bash
make install-tools
```

This installs:
- Air (hot reload)
- Swag (Swagger generator)
- golangci-lint (linter)
- goimports (formatter)
- migrate (database migrations)

### 3. Run with Docker (Recommended)

```bash
# Start development environment
make docker-dev

# View logs
make docker-logs
```

Access:
- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- pgAdmin: http://localhost:5050 (admin@admin.com / admin)
- Redis Commander: http://localhost:8081

### 4. Run Locally

```bash
# Download dependencies
make deps

# Run database migrations
make migrate-up

# Run with hot reload
make dev
```

## 📋 Makefile Commands

### Development
| Command | Description |
|---------|-------------|
| `make dev` | Run with Air hot reload |
| `make run` | Build and run |
| `make build` | Build binary |
| `make clean` | Clean build artifacts |

### Testing
| Command | Description |
|---------|-------------|
| `make test` | Run tests |
| `make test-coverage` | Run tests with coverage |
| `make test-verbose` | Run tests verbosely |

### Code Quality
| Command | Description |
|---------|-------------|
| `make lint` | Run linter |
| `make fmt` | Format code |
| `make tidy` | Tidy go.mod |
| `make vet` | Run go vet |

### Database
| Command | Description |
|---------|-------------|
| `make migrate-up` | Run migrations up |
| `make migrate-down` | Run migrations down |
| `make migrate-create NAME=name` | Create new migration |
| `make seed` | Run database seeder |

### Docker
| Command | Description |
|---------|-------------|
| `make docker-dev` | Start dev containers |
| `make docker-dev-down` | Stop dev containers |
| `make docker-prod` | Start prod containers |
| `make docker-build` | Build Docker image |

### Documentation
| Command | Description |
|---------|-------------|
| `make swagger` | Generate Swagger docs |
| `make swagger-fmt` | Format Swagger comments |

## 🔐 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user

### Users (Protected)
- `GET /api/v1/users/me` - Get current user
- `GET /api/v1/users` - Get all users (paginated)
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Health
- `GET /health` - Health check
- `GET /ready` - Readiness check

## 🔧 Configuration

Environment variables (`.env`):

```env
# Application
APP_NAME=go-clean-architecture
APP_ENV=development
APP_PORT=8080
APP_DEBUG=true

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_clean_db
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

## 📝 Adding New Features

### 1. Create Entity
```go
// internal/entity/product.go
type Product struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"size:255"`
    Price float64
}
```

### 2. Create Repository Interface
```go
// internal/repository/product_repository.go
type ProductRepository interface {
    Create(ctx context.Context, product *entity.Product) error
    FindByID(ctx context.Context, id uint) (*entity.Product, error)
}
```

### 3. Implement Repository
```go
// internal/repository/product_repository_impl.go
type productRepository struct { db *gorm.DB }
// ... implement methods
```

### 4. Create UseCase
```go
// internal/usecase/product_usecase.go
type ProductUseCase interface {
    Create(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error)
}
```

### 5. Create Handler with Swagger Annotations
```go
// internal/handler/product_handler.go
// @Summary Create product
// @Tags Products
// @Router /api/v1/products [post]
func (h *ProductHandler) Create(c *gin.Context) {}
```

### 6. Register Routes
```go
// internal/router/router.go
products := v1.Group("/products")
products.POST("", productHandler.Create)
```

### 7. Generate Swagger
```bash
make swagger
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Open coverage report
open coverage.html
```

## 🚢 Deployment

### Docker Production
```bash
# Build production image
make docker-build

# Start production containers
make docker-prod
```

### Manual Deployment
```bash
# Build binary
make build

# Run migrations
./bin/main migrate up

# Start server
./bin/main
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open Pull Request
