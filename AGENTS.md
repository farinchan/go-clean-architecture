# AGENTS.md — AI Agent Guidelines & Architecture Manual

> **Scope**: This document provides foundational context, architectural rules, code conventions, and step-by-step instructions for AI Agents (Antigravity, Claude, Copilot, Cursor, etc.) working on this repository.
> **Repository**: `go-clean-architecture`
> **Language & Runtime**: Go (Golang 1.21+, Go Toolchain 1.26+)

---

## 1. Project Overview & Tech Stack

This project is a production-grade, highly structured RESTful API template implementing **Uncle Bob's Clean Architecture** in Go.

### Core Stack
- **Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin) (`v1.12.0+`)
- **ORM**: [GORM](https://gorm.io/) (`v1.31.2+`) with PostgreSQL Driver (`gorm.io/driver/postgres`)
- **Database**: PostgreSQL 16
- **Cache**: Redis 7 via [go-redis/redis/v8](https://github.com/go-redis/redis)
- **Migrations**: [golang-migrate/migrate/v4](https://github.com/golang-migrate/migrate)
- **Auth**: JWT via [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) (HS256) & bcrypt password hashing
- **Validation**: [go-playground/validator/v10](https://github.com/go-playground/validator)
- **Configuration**: [Viper](https://github.com/spf13/viper) (reads `.env` and environment variables)
- **Logging**: [Logrus](https://github.com/sirupsen/logrus) (Structured JSON for production, colored text for dev)
- **Email**: [Gomail v2](https://gopkg.in/gomail.v2) for SMTP
- **API Documentation**: OpenAPI / Swagger via [swaggo/swag](https://github.com/swaggo/swag) and [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)
- **Hot Reload**: [Air](https://github.com/cosmtrek/air) (`.air.toml`)
- **Containerization**: Docker & Docker Compose (`docker-compose.yml`, `Dockerfile`)

---

## 2. Clean Architecture Layering & Invariant Rules

The application strictly obeys the **Dependency Inversion Principle**. Outer layers know about inner layers; **inner layers NEVER know about outer layers**.

```
[ HTTP Request ]
       ↓
[ internal/router & middleware ]
       ↓
[ internal/handler (HTTP Controllers) ]
       ↓
[ internal/usecase (Business Rules) ]  ← [ internal/dto (Data Contracts) ]
       ↓
[ internal/repository (Data Access Interface) ]
       ↓
[ internal/repository_impl (GORM) ]
       ↓
[ internal/entity (Database Models) ]  ← [ PostgreSQL / Redis ]
```

### Layer Boundaries & Restrictions:

1. **`internal/entity/` (Domain Entities)**
   - Pure Go structs representing persistent database records and GORM schema definitions.
   - **RULE**: NEVER import `github.com/gin-gonic/gin` or HTTP transport packages here.
   - Sensitive fields (e.g. `Password`) must have `json:"-"` and `gorm:"not null"`.
   - Soft deletes use `gorm.DeletedAt` with `gorm:"index"`.

2. **`internal/dto/` (Data Transfer Objects)**
   - Request and response structs for HTTP boundaries.
   - Validation rules must be defined via `binding:"..."` struct tags (e.g., `binding:"required,email"`).
   - Responses must never leak raw entities containing internal or sensitive fields.

3. **`internal/repository/` (Data Access Layer)**
   - Separated into **Interface** (`user_repository.go`) and **Implementation** (`user_repository_impl.go`).
   - Every repository method **MUST** accept `ctx context.Context` as the first parameter.
   - Use `r.db.WithContext(ctx)` for all GORM queries.
   - **RULE**: Business logic must NEVER be placed in the repository.

4. **`internal/usecase/` (Application Business Logic)**
   - Defined as an **Interface** (`UserUseCase`) and its concrete struct (`userUseCase`).
   - Accepts DTOs / primitives from handlers and interacts strictly with `Repository` interfaces.
   - Handles password hashing, password comparisons, status checks, and token generations.
   - **RULE**: NEVER pass `*gin.Context` into use cases; pass standard `context.Context`.

5. **`internal/handler/` (HTTP Presentation Layer)**
   - Parses parameters, headers, and request bodies (`c.ShouldBindJSON`).
   - Formats validation errors using `pkg/validator.FormatValidationErrors(err)`.
   - Sends standardized responses using `pkg/response`.
   - Must contain Swaggo annotations (`@Summary`, `@Tags`, `@Router`, etc.) above each handler func.

6. **`internal/router/` & `internal/middleware/`**
   - Configures Gin route groups, global middlewares (`RecoveryMiddleware`, `LoggerMiddleware`, `CORSMiddleware`), and authenticated routes (`AuthMiddleware`, `RoleMiddleware`).

7. **`pkg/` (Shared Infrastructure & Utilities)**
   - Reusable modules across the application:
     - `pkg/response`: Standard API response envelope (`Success`, `Created`, `BadRequest`, `NotFound`, `ValidationError`, etc.).
     - `pkg/validator`: Custom validator engine and error translator.
     - `pkg/utils`: `JWTManager`, `HashPassword`, `CheckPassword`, `GenerateUUID`, `GenerateOTP`, `SlugFromString`.
     - `pkg/database`: GORM PostgreSQL connection pool and Redis client wrapper.
     - `pkg/logger`: Logrus logging helpers.
     - `pkg/mail`: SMTP email sender.
     - `pkg/apperrors`: Domain error wrappers and status codes.

---

## 3. Standard Response Format

All HTTP endpoints must respond using [`pkg/response`](pkg/response/response.go). Never return raw `c.JSON()` maps directly.

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "is_active": true,
    "created_at": "2026-10-20T16:00:00Z",
    "updated_at": "2026-10-20T16:00:00Z"
  },
  "error": null,
  "meta": null
}
```

For paginated responses:
```go
meta := response.BuildMeta(page, limit, total)
response.SuccessWithMeta(c, "Data retrieved successfully", items, meta)
```

For validation errors (HTTP 422):
```go
errors := validator.FormatValidationErrors(err)
response.ValidationError(c, errors)
```

---

## 4. Step-by-Step: Adding a New Feature (e.g. "Product")

When implementing a new entity or resource, follow this exact sequence:

### Step 1: Database Migration
Create up & down SQL files in `database/migrations/`:
```bash
make migrate-create NAME=create_products_table
```
Write standard PostgreSQL DDL in the `.up.sql` and corresponding drop statements in `.down.sql`.

### Step 2: Entity (`internal/entity/product.go`)
```go
package entity

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"size:255;not null"`
    Price       float64        `json:"price" gorm:"not null"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Product) TableName() string {
    return "products"
}
```

### Step 3: DTO (`internal/dto/product_dto.go`)
```go
package dto

type CreateProductRequest struct {
    Name  string  `json:"name" binding:"required,min=2,max=255"`
    Price float64 `json:"price" binding:"required,gt=0"`
}

type ProductResponse struct {
    ID    uint    `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

### Step 4: Repository (`internal/repository/product_repository.go` & `product_repository_impl.go`)
Define the interface:
```go
type ProductRepository interface {
    Create(ctx context.Context, product *entity.Product) error
    FindByID(ctx context.Context, id uint) (*entity.Product, error)
    FindAll(ctx context.Context, page, limit int) ([]entity.Product, int64, error)
    Update(ctx context.Context, product *entity.Product) error
    Delete(ctx context.Context, id uint) error
}
```
Implement with `*gorm.DB` in `product_repository_impl.go`.

### Step 5: UseCase (`internal/usecase/product_usecase.go`)
Define business logic interface and implementation:
```go
type ProductUseCase interface {
    Create(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error)
    GetByID(ctx context.Context, id uint) (*dto.ProductResponse, error)
}
```

### Step 6: Handler (`internal/handler/product_handler.go`)
Add Swagger documentation comments above each method:
```go
// Create godoc
// @Summary Create product
// @Tags Products
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Create Product"
// @Security BearerAuth
// @Success 201 {object} response.Response{data=dto.ProductResponse}
// @Router /api/v1/products [post]
func (h *ProductHandler) Create(c *gin.Context) { ... }
```

### Step 7: Router (`internal/router/router.go`)
Register routes inside `SetupRoutes()`:
```go
products := v1.Group("/products")
products.Use(middleware.AuthMiddleware(r.jwtManager))
{
    products.POST("", r.productHandler.Create)
    products.GET("/:id", r.productHandler.GetByID)
}
```

### Step 8: Wire in `cmd/api/main.go`
Instantiate Repository $\rightarrow$ UseCase $\rightarrow$ Handler $\rightarrow$ Router.

### Step 9: Re-generate Swagger Docs & Verify
```bash
make swagger
go build -v ./cmd/api
```

---

## 5. Development Commands Reference

| Task | Command | Notes |
| :--- | :--- | :--- |
| **Hot Reload** | `make dev` | Uses `air` config in `.air.toml` |
| **Run Directly** | `make run` | Builds to `bin/main.exe` and executes |
| **Build Binary** | `make build` | Compiles `cmd/api/main.go` |
| **Clean Artifacts** | `make clean` | Removes `bin/`, `tmp/`, coverage files |
| **Run Tests** | `make test` | Executes `go test ./... -race` |
| **Test Coverage** | `make test-coverage`| Generates `coverage.html` |
| **Static Analysis** | `make vet` | Executes `go vet ./...` |
| **Format Code** | `make fmt` | Runs `go fmt` and `goimports -w .` |
| **Tidy Modules** | `make tidy` | Runs `go mod tidy` |
| **Run Migrations** | `make migrate-up` | Runs all `.up.sql` migrations |
| **Rollback Migration**| `make migrate-down` | Rolls back migrations |
| **Migration Version**| `make migrate-version`| Prints current migration version |
| **Database Seeding**| `make seed` | Runs `cmd/seed/main.go` |
| **Generate Swagger**| `make swagger` | Runs `swag init` for `docs/` |
| **Start Docker** | `make docker-up` | Starts app, PostgreSQL, and Redis |
| **Stop Docker** | `make docker-down` | Stops containers |
| **Docker Logs** | `make docker-logs` | Follows container logs |
| **Docker Clean** | `make docker-clean`| Removes containers & persistent volumes |

---

## 6. Crucial Agent Guardrails (Do's and Don'ts)

### ❌ What Agents MUST NEVER Do:
1. **DO NOT bypass the UseCase or Repository**: Never perform GORM DB queries directly from a Gin Handler.
2. **DO NOT pass Gin context deeper than Handler**: Never pass `*gin.Context` to UseCases or Repositories; always pass `c.Request.Context()` as `context.Context`.
3. **DO NOT edit generated Swagger files directly**: Never hand-edit `docs/docs.go`, `docs/swagger.json`, or `docs/swagger.yaml`. Edit annotations in `internal/handler/` or `cmd/api/main.go` and run `make swagger`.
4. **DO NOT commit `.env` files**: Secret keys and DB credentials belong in `.env` (gitignored). Keep `.env.example` updated with mock defaults.
5. **DO NOT mix transport concerns into entities**: Entities must not have `form` or `binding` tags; those belong strictly in `internal/dto/`.
6. **DO NOT remove existing comments or docstrings** unrelated to your specific code edits.

### ✅ What Agents MUST ALWAYS Do:
1. Always run `go mod tidy` and `go vet ./...` after adding or changing dependencies.
2. Always ensure all 3 entry points (`cmd/api`, `cmd/migrate`, `cmd/seed`) compile cleanly:
   ```bash
   go build ./cmd/api ./cmd/migrate ./cmd/seed
   ```
3. Use `pkg/response` functions for returning HTTP responses to maintain uniform API envelopes.
4. Keep interfaces in the target layer and instantiate them via constructor functions (e.g., `NewUserRepository`, `NewUserUseCase`).
5. Ensure database migrations are idempotent where applicable (`CREATE TABLE IF NOT EXISTS`, `CREATE INDEX IF NOT EXISTS`).
