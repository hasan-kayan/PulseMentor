# Yeni Sistem Ekleme Rehberi

Bu rehber, PulseMentor backend projesine yeni bir domain/sistem eklerken izlemeniz gereken adımları açıklar.

## Proje Yapısı

Proje **Domain-Driven Design (DDD)** prensiplerine göre organize edilmiştir:

```
Backend/
├── cmd/server/          # Uygulama giriş noktası
├── internal/
│   ├── app/            # Uygulama başlatma ve bağımlılık yönetimi
│   ├── config/         # Konfigürasyon yönetimi
│   ├── domain/         # Domain modelleri ve business logic
│   │   ├── users/      # Örnek: Users domain
│   │   └── [yeni_domain]/
│   ├── http/           # HTTP katmanı
│   │   ├── handlers/   # HTTP handler'lar
│   │   ├── middleware/ # Middleware'ler
│   │   └── routes/     # Route tanımlamaları
│   ├── infra/          # Altyapı katmanı
│   │   ├── db/postgres/ # Database implementasyonları
│   │   └── ...
│   └── shared/         # Paylaşılan yardımcı fonksiyonlar
├── migrations/         # Database migration dosyaları
└── go.mod
```

## Yeni Bir Domain Ekleme Adımları

### 1. Domain Model Oluşturma

`internal/domain/[domain_adi]/` klasörü altında domain dosyalarını oluşturun:

#### `model.go` - Domain modelleri
```go
package [domain_adi]

import "time"

type [Entity] struct {
    ID        string
    // Diğer alanlar
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Create[Entity]Input struct {
    // Gerekli alanlar
}
```

#### `repository.go` - Repository interface
```go
package [domain_adi]

import "context"

type Repository interface {
    Create(ctx context.Context, entity *[Entity]) error
    FindByID(ctx context.Context, id string) (*[Entity], error)
    // Diğer metodlar
    Update(ctx context.Context, entity *[Entity]) error
    Delete(ctx context.Context, id string) error
}
```

#### `service.go` - Business logic
```go
package [domain_adi]

import (
    "context"
    "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/config"
    sharedErrors "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/errors"
    "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/id"
    "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/validate"
)

type Service struct {
    repo   Repository
    config config.Config
}

func NewService(repo Repository, cfg config.Config) *Service {
    return &Service{
        repo:   repo,
        config: cfg,
    }
}

func (s *Service) Create(ctx context.Context, input Create[Entity]Input) (*[Entity], error) {
    // Validation
    // Business logic
    // Repository çağrısı
}
```

### 2. Database Repository Implementasyonu

`internal/infra/db/postgres/[domain_adi]_repo.go` dosyasını oluşturun:

```go
package postgres

import (
    "context"
    "fmt"
    "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/domain/[domain_adi]"
    "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/errors"
    "github.com/jackc/pgx/v5"
)

type [Entity]Repository struct {
    db *DB
}

func New[Entity]Repository(db *DB) *[Entity]Repository {
    return &[Entity]Repository{db: db}
}

func (r *[Entity]Repository) Create(ctx context.Context, entity *[domain_adi].[Entity]) error {
    query := `
        INSERT INTO [table_name] (id, ...)
        VALUES ($1, ...)
    `
    _, err := r.db.Pool.Exec(ctx, query, entity.ID, ...)
    if err != nil {
        return fmt.Errorf("failed to create [entity]: %w", err)
    }
    return nil
}

// Diğer metodlar...
```

### 3. Database Migration Oluşturma

`migrations/` klasörüne yeni migration dosyası ekleyin:

```sql
-- migrations/XXXX_create_[table_name].sql
CREATE TABLE IF NOT EXISTS [table_name] (
    id VARCHAR(32) PRIMARY KEY,
    -- Diğer kolonlar
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Gerekli index'ler
CREATE INDEX IF NOT EXISTS idx_[table_name]_[column] ON [table_name]([column]);
```

### 4. HTTP Handler Oluşturma

`internal/http/handlers/[domain_adi]_handler.go` dosyasını oluşturun:

```go
package handlers

import (
    "net/http"
    "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/domain/[domain_adi]"
    "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/http/httpx"
    "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/errors"
)

type [Entity]Handler struct {
    [entity]Service *[domain_adi].Service
}

func New[Entity]Handler([entity]Service *[domain_adi].Service) *[Entity]Handler {
    return &[Entity]Handler{
        [entity]Service: [entity]Service,
    }
}

type Create[Entity]Request struct {
    // Request alanları
}

func (h *[Entity]Handler) Create(w http.ResponseWriter, r *http.Request) {
    var req Create[Entity]Request
    if err := httpx.BindJSON(r, &req); err != nil {
        httpx.Error(w, http.StatusBadRequest, err)
        return
    }

    entity, err := h.[entity]Service.Create(r.Context(), [domain_adi].Create[Entity]Input{
        // Mapping
    })
    if err != nil {
        status := http.StatusInternalServerError
        if err == errors.ErrInvalidInput || err == errors.ErrAlreadyExists {
            status = http.StatusBadRequest
        }
        httpx.Error(w, status, err)
        return
    }

    httpx.JSON(w, http.StatusCreated, entity)
}
```

### 5. Route Tanımlama

`internal/http/routes/routes.go` dosyasını güncelleyin:

```go
func SetupRouter(userService *users.Service, [entity]Service *[domain_adi].Service) *chi.Mux {
    // ...
    
    // Yeni handler'ı oluştur
    [entity]Handler := handlers.New[Entity]Handler([entity]Service)
    
    // Route'ları ekle
    r.Route("/api/v1", func(r chi.Router) {
        // Public routes
        // ...
        
        // Protected routes
        r.Group(func(r chi.Router) {
            r.Use(authMiddleware.RequireAuth)
            r.Post("/[entities]", [entity]Handler.Create)
            r.Get("/[entities]", [entity]Handler.List)
            r.Get("/[entities]/{id}", [entity]Handler.Get)
            r.Put("/[entities]/{id}", [entity]Handler.Update)
            r.Delete("/[entities]/{id}", [entity]Handler.Delete)
        })
    })
}
```

### 6. App Başlatma

`internal/app/app.go` dosyasını güncelleyin:

```go
func New(cfg config.Config) (*App, error) {
    // ...
    
    // Yeni repository'yi oluştur
    [entity]Repo := postgres.New[Entity]Repository(db)
    
    // Yeni service'i oluştur
    [entity]Service := [domain_adi].NewService([entity]Repo, cfg)
    
    // Router'a service'i ekle
    router := routes.SetupRouter(userService, [entity]Service)
    
    // ...
}
```

## Authentication Kullanımı

### Protected Route Oluşturma

Kullanıcı kimlik doğrulaması gerektiren route'lar için `authMiddleware.RequireAuth` kullanın:

```go
r.Group(func(r chi.Router) {
    r.Use(authMiddleware.RequireAuth)
    r.Get("/protected", handler.ProtectedHandler)
})
```

### Handler'da Kullanıcı ID'sini Alma

Handler içinde authenticated kullanıcının ID'sini almak için:

```go
import "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/context"

func (h *Handler) SomeHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := context.GetUserID(r.Context())
    if !ok {
        httpx.Error(w, http.StatusUnauthorized, errors.ErrUnauthorized)
        return
    }
    
    // userID'yi kullan
}
```

## Örnek: "Posts" Domain'i Ekleme

Tam bir örnek için aşağıdaki adımları takip edin:

1. **Domain oluştur:**
   - `internal/domain/posts/model.go`
   - `internal/domain/posts/repository.go`
   - `internal/domain/posts/service.go`

2. **Repository implementasyonu:**
   - `internal/infra/db/postgres/posts_repo.go`

3. **Migration:**
   - `migrations/0002_create_posts.sql`

4. **Handler:**
   - `internal/http/handlers/posts_handler.go`

5. **Route ekle:**
   - `internal/http/routes/routes.go` dosyasını güncelle

6. **App'e bağla:**
   - `internal/app/app.go` dosyasını güncelle

## Best Practices

1. **Error Handling:** Her zaman `shared/errors` paketindeki standart hataları kullanın
2. **Validation:** `shared/validate` paketindeki validation fonksiyonlarını kullanın
3. **ID Generation:** `shared/id.New()` fonksiyonunu kullanın
4. **Context:** Her zaman context'i metodlara geçirin
5. **Password:** Asla password'ü response'da döndürmeyin
6. **Type Safety:** Context key'leri için `shared/context` paketini kullanın

## Yardımcı Paketler

### `shared/errors`
Standart hata tanımlamaları:
- `ErrNotFound`
- `ErrUnauthorized`
- `ErrForbidden`
- `ErrInvalidInput`
- `ErrAlreadyExists`
- `ErrInternal`

### `shared/validate`
Validation fonksiyonları:
- `Email(email string) bool`
- `Password(password string) bool`
- `NonEmpty(s string) bool`

### `shared/id`
ID generation:
- `New() string` - 32 karakterlik hex ID üretir

### `shared/context`
Context helper'ları:
- `GetUserID(ctx context.Context) (string, bool)`
- `WithUserID(ctx context.Context, userID string) context.Context`

### `httpx`
HTTP helper'ları:
- `BindJSON(r *http.Request, v interface{}) error`
- `JSON(w http.ResponseWriter, status int, data interface{})`
- `Error(w http.ResponseWriter, status int, err error)`

## Environment Variables

Yeni bir domain eklerken gerekirse `internal/config/config.go` ve `internal/config/load.go` dosyalarına yeni konfigürasyon değerleri ekleyin.

## Test Etme

1. Migration'ı çalıştırın
2. Server'ı başlatın
3. API endpoint'lerini test edin (Postman, curl, vs.)
4. Authentication gerektiren endpoint'leri test ederken `Authorization: Bearer <token>` header'ını ekleyin

## Sorun Giderme

- **Import hatası:** `go mod tidy` çalıştırın
- **Database bağlantı hatası:** `DATABASE_URL` environment variable'ını kontrol edin
- **JWT hatası:** `JWT_SECRET` environment variable'ını kontrol edin
- **Migration hatası:** Migration dosyalarının sırasını kontrol edin

## Özet Checklist

- [ ] Domain model dosyaları oluşturuldu (`model.go`, `repository.go`, `service.go`)
- [ ] Database repository implementasyonu yapıldı
- [ ] Migration dosyası oluşturuldu ve çalıştırıldı
- [ ] HTTP handler oluşturuldu
- [ ] Route'lar tanımlandı
- [ ] App başlatma koduna eklendi
- [ ] Test edildi

Bu rehberi takip ederek yeni domain'leri tutarlı ve bakımı kolay bir şekilde ekleyebilirsiniz.

