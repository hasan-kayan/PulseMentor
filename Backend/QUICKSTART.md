# Hızlı Başlangıç - PulseMentor Backend

## Adım 1: PostgreSQL Veritabanı Oluştur

```bash
# PostgreSQL'e bağlan
psql -U postgres

# Veritabanı oluştur
CREATE DATABASE pulsementor;

# Çık
\q
```

## Adım 2: Environment Variables Ayarla

Terminal'de çalıştır:

```bash
export DATABASE_URL="postgres://postgres@localhost:5432/pulsementor?sslmode=disable"
export JWT_SECRET="dev-secret-key-change-in-production-min-32-chars-long"
```

**Not:** `postgres` kullanıcı adını kendi PostgreSQL kullanıcı adınızla değiştirin. Şifre varsa:
```bash
export DATABASE_URL="postgres://username:password@localhost:5432/pulsementor?sslmode=disable"
```

## Adım 3: Migration'ları Çalıştır

```bash
cd Backend
psql $DATABASE_URL -f migrations/0001_create_users.sql
```

## Adım 4: Server'ı Başlat

### Seçenek 1: Development Script (Önerilen)

```bash
./scripts/dev.sh
```

### Seçenek 2: Manuel

```bash
go run cmd/server/main.go
```

### Seçenek 3: Build edip çalıştır

```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

## Server Başarıyla Başladığında

Şu mesajı göreceksiniz:
```
listening on 0.0.0.0:8080
```

## Test Et

Yeni bir terminal açıp:

```bash
# Kullanıcı kaydı
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Giriş yap
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## Sorun Giderme

### "DATABASE_URL is required" hatası
- Environment variable'ı ayarladığınızdan emin olun
- Terminal'de `echo $DATABASE_URL` ile kontrol edin

### "connection refused" hatası
- PostgreSQL'in çalıştığından emin olun: `pg_isready`
- DATABASE_URL'deki host ve port'u kontrol edin

### "database does not exist" hatası
- Adım 1'deki veritabanı oluşturma adımını tekrar kontrol edin

### "JWT_SECRET is required" hatası
- Environment variable'ı ayarladığınızdan emin olun

## Tüm Environment Variables (Opsiyonel)

```bash
export APP_ENV=dev
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
export DATABASE_URL="postgres://postgres@localhost:5432/pulsementor?sslmode=disable"
export JWT_SECRET="dev-secret-key-change-in-production-min-32-chars-long"
export JWT_ISSUER="pulsementor"
export ACCESS_TOKEN_TTL="24h"
export REFRESH_TOKEN_TTL="168h"
export BCRYPT_COST="12"
```

