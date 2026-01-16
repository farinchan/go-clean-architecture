# Go Clean Architecture Backend

Clean Architecture Go Backend dengan Gin, GORM, PostgreSQL, dan Redis.

## Struktur Project

```
go-clean-architecture/
├── cmd/                        # Entry points
│   └── api/
│       └── main.go
├── config/                     # Konfigurasi aplikasi
│   └── config.go
├── infrastructure/             # Implementasi infrastructure
│   ├── cache/
│   │   └── redis.go
│   └── database/
│       └── postgres.go
├── internal/                   # Business logic
│   ├── domain/                 # Entities & Repository interfaces
│   │   ├── user.go
│   │   └── user_repository.go
│   ├── dto/                    # Data Transfer Objects
│   │   └── user_dto.go
│   ├── repository/             # Repository implementations
│   │   └── user_repository.go
│   ├── usecase/                # Business logic / Use cases
│   │   └── user_usecase.go
│   └── delivery/               # HTTP handlers & routes
│       └── http/
│           ├── handler/
│           │   └── user_handler.go
│           ├── middleware/
│           │   └── middleware.go
│           └── route/
│               └── route.go
├── migrations/                 # Database migrations
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
├── pkg/                        # Shared packages
│   ├── logger/
│   │   └── logger.go
│   ├── response/
│   │   └── response.go
│   └── validator/
│       └── validator.go
├── .env                        # Environment configuration
├── .env.example                # Environment template
├── go.mod
├── go.sum
└── README.md
```

## Teknologi

- **Gin** - HTTP Web Framework
- **GORM** - ORM untuk database
- **Golang Migrate** - Database migrations
- **Go Playground Validator** - Validasi request
- **Viper** - Configuration management
- **Logrus** - Logging
- **PostgreSQL** - Database
- **Redis** - Caching

## Instalasi

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Setup Database PostgreSQL

Buat database:
```sql
CREATE DATABASE go_clean_architecture_db;
```

### 3. Install Golang Migrate CLI

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 4. Jalankan Migration

```bash
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/go_clean_architecture_db?sslmode=disable" up
```

### 5. Konfigurasi

Copy file `.env.example` ke `.env` dan sesuaikan:

```bash
cp .env.example .env
```

Edit file `.env` sesuai environment:

```env
# Application
APP_NAME=go-clean-architecture
APP_PORT=8080
APP_ENV=development

# Database PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_clean_architecture_db
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 6. Jalankan Aplikasi

```bash
go run cmd/api/main.go
```

### 7. Build Aplikasi

```bash
go build -o bin/api cmd/api/main.go
./bin/api
```

## API Endpoints

### Health Check
- `GET /health` - Check server status

### Users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - Get all users (dengan pagination)
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

## Contoh Request

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "password123"}'
```

### Get All Users
```bash
curl http://localhost:8080/api/v1/users?page=1&limit=10
```

## Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`)
   - Berisi entities dan repository interfaces
   - Tidak bergantung pada layer lain

2. **Use Case Layer** (`internal/usecase/`)
   - Berisi business logic
   - Bergantung pada domain layer

3. **Repository Layer** (`internal/repository/`)
   - Implementasi akses data
   - Bergantung pada domain layer

4. **Delivery Layer** (`internal/delivery/`)
   - HTTP handlers dan routes
   - Bergantung pada use case layer
