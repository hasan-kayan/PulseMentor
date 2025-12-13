# PulseMentor Backend - Kurulum ve Çalıştırma

## Gereksinimler

- Go 1.21 veya üzeri
- PostgreSQL 12 veya üzeri
- psql (PostgreSQL client)

## Hızlı Başlangıç

### 1. Bağımlılıkları Yükle

```bash
cd Backend
go mod download
```

### 2. PostgreSQL Veritabanı Oluştur

```bash
createdb pulsementor
```

veya psql ile:

```bash
psql -U postgres
CREATE DATABASE pulsementor;
\q
```

### 3. Environment Variables Ayarla

Aşağıdaki environment variable'ları ayarlayın:

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/pulsementor?sslmode=disable"
export JWT_SECRET="your-super-secret-jwt-key-min-32-chars"
export JWT_ISSUER="pulsementor"
export ACCESS_TOKEN_TTL="24h"
export REFRESH_TOKEN_TTL="168h"
export BCRYPT_COST="12"
export SERVER_HOST="0.0.0.0"
export SERVER_PORT="8080"
```

**Not:** Production'da mutlaka güçlü bir `JWT_SECRET` kullanın (en az 32 karakter).

### 4. Migration'ları Çalıştır

```bash
./scripts/migrate.sh
```

veya manuel olarak:

```bash
psql $DATABASE_URL -f migrations/0001_create_users.sql
```

### 5. Server'ı Başlat

#### Development Mode (Script ile):

```bash
./scripts/dev.sh
```

#### Manuel:

```bash
go run cmd/server/main.go
```

#### Build edip çalıştırma:

```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

## API Endpoints

Server başladıktan sonra aşağıdaki endpoint'ler kullanılabilir:

### Public Endpoints

- `POST /api/v1/auth/register` - Kullanıcı kaydı
- `POST /api/v1/auth/login` - Giriş yapma
- `POST /api/v1/auth/refresh` - Token yenileme

### Protected Endpoints (Authorization: Bearer <token> gerekli)

- `GET /api/v1/auth/me` - Kullanıcı bilgileri

## Örnek Kullanım

### Kullanıcı Kaydı

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Giriş Yapma

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "token": {
      "access_token": "...",
      "refresh_token": "..."
    }
  }
}
```

### Kullanıcı Bilgileri (Protected)

```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Token Yenileme

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

## Sorun Giderme

### Database Bağlantı Hatası

- `DATABASE_URL` environment variable'ının doğru olduğundan emin olun
- PostgreSQL'in çalıştığından emin olun: `pg_isready`
- Veritabanının var olduğundan emin olun

### JWT Secret Hatası

- `JWT_SECRET` environment variable'ının ayarlandığından emin olun
- Production'da güçlü bir secret kullanın

### Migration Hatası

- Migration dosyalarının sırasını kontrol edin
- Veritabanı bağlantısını kontrol edin
- Migration'ların daha önce çalıştırılmadığından emin olun

## Development

### Kod Formatlama

```bash
go fmt ./...
```

### Linting

```bash
go vet ./...
```

### Test

```bash
go test ./...
```

## Production Deployment

Production'da:

1. Güçlü `JWT_SECRET` kullanın
2. `APP_ENV=production` ayarlayın
3. SSL/TLS kullanın
4. Rate limiting ekleyin
5. Monitoring ve logging kurun
6. Database backup stratejisi oluşturun

