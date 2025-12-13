# AI Backend (Go) — Proje Yapısı Mantığı ve Çalışma Şekli (Monolith)

Bu döküman, oluşturduğun klasör yapısının:
- **neden böyle tasarlandığını**
- **request geldiğinde sistemin nasıl çalıştığını**
- her klasör/dosya tipinin **ne işe yaradığını**
basit ve net şekilde açıklar.

---

## 1) Büyük Resim: Monolith ama Modüler

Bu yapı **tek deploy edilen** bir backend (monolith) üretir:
- Tek binary → tek container → tek release.
Ama içeride modülerdir:
- **Domain** (iş kuralları) ayrı,
- **Infra** (DB/LLM/Redis gibi dış dünya) ayrı,
- **HTTP** (endpoint/middleware) ayrı.

Bu sayede:
- Kod büyüdükçe dağılmaz,
- Test yazmak kolaylaşır,
- İleride mikroservise bölmek istersen “koparması” rahat olur.

---

## 2) Request Akışı: Bir İstek Sistemde Nasıl Gezer?

Örnek: `POST /chat`

1. **`cmd/server/main.go`** çalışır  
   → uygulamayı ayağa kaldırır.

2. **`internal/app`** başlatır  
   - config okur (`internal/config`)
   - logger kurar (`internal/observability`)
   - DB bağlantısını açar (`internal/infra/db/postgres`)
   - Redis’i açar (`internal/infra/cache/redis`)
   - LLM client’ı kurar (`internal/infra/llm/openai`)
   - router’ı oluşturur (`internal/http/routes`)

3. Request server’a gelir → **middleware zinciri** çalışır (`internal/http/middleware`)
   - request id
   - logging
   - auth
   - rate limit
   - recover

4. Route eşleşir → **handler** çağrılır (`internal/http/handlers`)
   - request body parse eder
   - validate eder
   - domain service çağırır

5. **domain service** iş kuralını yürütür (`internal/domain/.../service.go`)
   - gerekiyorsa repo çağırır (DB)
   - gerekiyorsa LLM çağırır
   - sonucu “domain modeli” olarak üretir

6. Repo implementasyonu **infra** içindedir (`internal/infra/...`)
   - DB sorguları
   - Redis cache
   - LLM API çağrıları

7. Handler sonucu JSON’a çevirir ve response döner.

---

## 3) “Domain vs Infra” Mantığı (En Kritik Kısım)

### Domain (iş kuralı)
- “Ne yapıyoruz?” sorusunun cevabı.
- Örn: “Kullanıcının günlük token limitini aşmadığından emin ol, prompt’u hazırla, konuşmayı kaydet.”

### Infra (teknik entegrasyon)
- “Nasıl yapıyoruz?” sorusunun cevabı.
- Örn: “Postgres’e nasıl bağlanıyorum?”, “OpenAI API’ye hangi URL ile gidiyorum?”

**Kural:**
- Domain, OpenAI paketini bilmez.
- Domain, sadece `LLMClient` gibi bir interface bilir.
- OpenAI implementasyonu infra’da durur.

---

## 4) Klasörlerin Mantığı ve Ne İşe Yaradıkları

## `cmd/server/`
**Çalıştırma girişi (entrypoint).**
- Uygulama nereden başlar? → buradan.
- İçine ne koyarsın:
  - `main.go` (minimum kod, sadece başlatma)

**Mantık:**  
“Server’ı koşacağım yer burası, iş kuralı burada olmaz.”

---

## `internal/app/`
**Uygulamanın kurulum ve bağlama katmanı (wiring).**
- Tüm bağımlılıkları bir araya getirir:
  - config + logger + db + redis + llm + router
- Graceful shutdown gibi yaşam döngüsü yönetimi burada olur.

**İçine ne koyarsın:**
- `app.go` (New/Run)
- `deps.go` (db/redis/llm init)
- `shutdown.go`

**Mantık:**  
“Her şeyi tek yerde kur, geri kalan katmanlar birbirini direkt tanımasın.”

---

## `internal/config/`
**Config yükleme ve doğrulama.**
- `.env`, environment variables, defaults, validation.

**İçine ne koyarsın:**
- `config.go` (Config struct)
- `load.go` (env read)
- `validate.go`

**Mantık:**  
“Uygulama ayarları tek noktadan yönetilsin. Prod/Dev farkı burada çözülür.”

---

## `internal/http/`
**HTTP API katmanı (dış dünyaya açılan kapı).**

### `internal/http/middleware/`
Her request’i saran ortak işler.

**Ne işe yarar:**
- Güvenlik (auth)
- Stabilite (panic recover)
- Limit (rate limit)
- İzlenebilirlik (log, request id)

**Mantık:**  
“Handler’lar temiz kalsın; ortak işler her endpoint’e tekrar yazılmasın.”

### `internal/http/routes/`
Endpoint kayıtları.

**Ne işe yarar:**
- Tüm route map tek yerde.
- Versiyonlama kolay: `/v1/...`

**Mantık:**  
“Uygulamanın dış kontratı (API) burada okunur.”

### `internal/http/handlers/`
HTTP handler’lar.

**Ne işe yarar:**
- JSON parse/serialize
- status code set
- domain service çağırma

**Mantık:**  
“Handler, sadece HTTP konuşur. İş kuralı yazmaz.”

---

## `internal/domain/`
**İş kurallarının merkezi.**

### `internal/domain/chat/`
- konuşma akışı
- mesaj formatı
- tool-calling kararları (varsa)
- konuşma kaydı

**Dosya tipleri:**
- `model.go`: domain entity’ler
- `service.go`: iş akışı
- `repository.go`: DB’ye ihtiyaç varsa interface

### `internal/domain/users/`
- register/login
- rol/izin
- session yönetimi (mantık)

### `internal/domain/usage/`
- token/kredi hesaplama
- quota kontrol
- usage record üretimi

**Mantık:**  
“Projenin gerçek değeri burada. Bu katman değişince ürün değişir.”

---

## `internal/infra/`
**Dış dünya adaptörleri.**

### `internal/infra/db/postgres/`
- Postgres bağlantısı
- SQL sorguları
- domain repository implementasyonları

**Mantık:**  
“DB detayı burada kalsın. Domain sadece interface görsün.”

### `internal/infra/cache/redis/`
- rate limit state
- session cache
- short-lived veriler

**Mantık:**  
“Cache bir optimizasyon/altyapı; domain’i kirletmesin.”

### `internal/infra/llm/openai/`
- OpenAI API client
- retry/backoff
- mapping

**Mantık:**  
“Provider değişirse sadece burayı değiştir.”

---

## `internal/jobs/`
**Arka plan işleri (async).**
- embedding üretme
- summarization
- indexing
- file processing

**Mantık:**  
“Uzun süren işleri request içinde bloklama; kuyruğa at, worker yapsın.”

---

## `internal/observability/`
**Sistemi görünür yapan katman.**
- structured logging
- metrics
- tracing

**Mantık:**  
“Prod’da problem olduğunda ‘ne oldu’yu burada anlarsın.”

---

## `internal/shared/`
**Her yerde kullanılan küçük yardımcılar.**

### `errors/`
- typed errors
- error mapping

### `validate/`
- ortak validasyon helper’ları

### `id/`
- UUID/ULID üretimi

**Mantık:**  
“Tekrar eden küçük parçalar tek yerde dursun.”

---

## 5) Proje Kök Klasörleri Ne İşe Yarar?

### `migrations/`
DB schema değişiklikleri.
- Mantık: “DB versiyonlanabilir olmalı.”

### `configs/`
örnek env dosyaları.
- Mantık: “Yeni developer hızlı ayağa kalksın.”

### `scripts/`
dev yardımcı scriptleri.
- Mantık: “Tek komutla setup/test/migrate.”

### `docs/`
mimari ve kullanım dokümanları.
- Mantık: “Kod kadar doküman da düzenli dursun.”

### `test/integration/`
gerçek DB/Redis ile uçtan uca testler.
- Mantık: “Gerçek hayata yakın test.”

### `api/`
OpenAPI sözleşmesi.
- Mantık: “Frontend/Client ile kontrat net olsun.”

---

## 6) Bu Yapıda “Ne Nereye Yazılır?” Mini Rehber

- Yeni endpoint:  
  `internal/http/handlers/...` + `internal/http/routes/routes.go`

- Yeni iş kuralı:  
  `internal/domain/<feature>/service.go`

- DB tablosu/repo:  
  interface → `internal/domain/<feature>/repository.go`  
  implementasyon → `internal/infra/db/postgres/<feature>_repo.go`

- LLM provider değişikliği:  
  `internal/infra/llm/...`

- Rate limit / auth / logging:  
  `internal/http/middleware/`

- Arka plan işi:  
  `internal/jobs/`

---

## 7) Hızlı Örnek Senaryo (Chat)

- Handler: `chat_handler.go`  
  parse eder → `ChatService.SendMessage(...)` çağırır

- Service: `domain/chat/service.go`  
  quota kontrol (`usage`) → prompt hazırlar → LLM çağırır → konuşmayı kaydeder

- LLM çağrısı: `infra/llm/openai/client.go`  
- DB kayıt: `infra/db/postgres/chat_repo.go`

---

Bu yapı ile ilerlerken kural basit:
**HTTP konuşur → Domain karar verir → Infra uygular.**



# Strucutre

ai-backend-go/
├─ go.mod
├─ go.sum
├─ Makefile
├─ README.md
├─ .gitignore
├─ .env.example
│
├─ cmd/
│  └─ server/
│     └─ main.go
│
├─ internal/
│  ├─ app/
│  │  ├─ app.go
│  │  ├─ deps.go
│  │  ├─ shutdown.go
│  │  └─ health.go
│  │
│  ├─ config/
│  │  ├─ config.go
│  │  ├─ load.go
│  │  └─ validate.go
│  │
│  ├─ http/
│  │  ├─ handlers/
│  │  │  ├─ chat_handler.go
│  │  │  ├─ users_handler.go
│  │  │  ├─ usage_handler.go
│  │  │  └─ health_handler.go
│  │  │
│  │  ├─ middleware/
│  │  │  ├─ request_id.go
│  │  │  ├─ logging.go
│  │  │  ├─ recover.go
│  │  │  ├─ auth.go
│  │  │  ├─ rate_limit.go
│  │  │  └─ cors.go
│  │  │
│  │  ├─ routes/
│  │  │  ├─ routes.go
│  │  │  └─ v1.go
│  │  │
│  │  └─ httpx/
│  │     ├─ response.go
│  │     ├─ errors.go
│  │     └─ bind.go
│  │
│  ├─ domain/
│  │  ├─ chat/
│  │  │  ├─ model.go
│  │  │  ├─ service.go
│  │  │  ├─ repository.go
│  │  │  └─ ports.go
│  │  │
│  │  ├─ users/
│  │  │  ├─ model.go
│  │  │  ├─ service.go
│  │  │  ├─ repository.go
│  │  │  └─ ports.go
│  │  │
│  │  └─ usage/
│  │     ├─ model.go
│  │     ├─ service.go
│  │     ├─ repository.go
│  │     └─ ports.go
│  │
│  ├─ infra/
│  │  ├─ db/
│  │  │  └─ postgres/
│  │  │     ├─ db.go
│  │  │     ├─ tx.go
│  │  │     ├─ migrations.go
│  │  │     ├─ chat_repo.go
│  │  │     ├─ users_repo.go
│  │  │     └─ usage_repo.go
│  │  │
│  │  ├─ cache/
│  │  │  └─ redis/
│  │  │     ├─ client.go
│  │  │     ├─ ratelimit_store.go
│  │  │     └─ session_store.go
│  │  │
│  │  └─ llm/
│  │     └─ openai/
│  │        ├─ client.go
│  │        ├─ mapper.go
│  │        └─ retry.go
│  │
│  ├─ jobs/
│  │  ├─ worker.go
│  │  ├─ scheduler.go
│  │  ├─ embeddings_job.go
│  │  └─ summarizer_job.go
│  │
│  ├─ observability/
│  │  ├─ logger.go
│  │  ├─ metrics.go
│  │  └─ tracing.go
│  │
│  └─ shared/
│     ├─ errors/
│     │  ├─ errors.go
│     │  └─ http_mapping.go
│     ├─ validate/
│     │  ├─ validate.go
│     │  └─ rules.go
│     ├─ id/
│     │  └─ id.go
│     └─ timeutil/
│        └─ timeutil.go
│
├─ migrations/
│  ├─ 0001_init.sql
│  ├─ 0002_chat.sql
│  ├─ 0003_users.sql
│  └─ 0004_usage.sql
│
├─ api/
│  └─ openapi.yaml
│
├─ configs/
│  ├─ dev.env
│  └─ prod.env
│
├─ scripts/
│  ├─ dev.sh
│  ├─ test.sh
│  ├─ lint.sh
│  └─ migrate.sh
│
├─ docs/
│  ├─ architecture.md
│  ├─ runbook.md
│  └─ decisions/
│     └─ 0001_structure.md
│
└─ test/
   └─ integration/
      ├─ chat_test.go
      ├─ users_test.go
      └─ usage_test.go
