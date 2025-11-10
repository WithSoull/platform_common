# Messenger Platform Library
Общая платформенная библиотека для микросервисов [мессенджера](https://github.com/WithSoull/messenger-overview). Содержит переиспользуемые компоненты отказоустойчивости, observability, middleware и утилиты для работы с инфраструктурой.

## Технологический стек
<p>
<img src="https://github.com/devicons/devicon/blob/master/icons/go/go-original.svg" title="Go" **alt="Go" width="60" height="60"/>&nbsp;
<img src="https://github.com/devicons/devicon/blob/master/icons/grpc/grpc-original.svg" title="gRPC" **alt="gRPC" width="60" height="60"/>&nbsp;
<img src="https://github.com/devicons/devicon/blob/master/icons/opentelemetry/opentelemetry-original.svg" title="OTEL" **alt="OTEL" width="60" height="60"/>&nbsp;
<img src="https://github.com/devicons/devicon/blob/master/icons/prometheus/prometheus-original.svg" title="Prometheus" **alt="Prometheus" width="60" height="60"/>&nbsp;
<img src="https://github.com/devicons/devicon/blob/master/icons/grafana/grafana-original.svg" title="Grafana"  alt="Grafana" width="60" height="60"/>&nbsp;
<img src="https://github.com/devicons/devicon/blob/master/icons/elasticsearch/elasticsearch-original.svg" title="ELS" **alt="ELS" width="60" height="60"/>&nbsp;
<img src="https://github.com/devicons/devicon/blob/master/icons/kibana/kibana-original.svg" title="Kibana" **alt="Kibana" width="60" height="60"/>&nbsp;
<img src="https://github.com/devicons/devicon/blob/master/icons/jaegertracing/jaegertracing-original.svg" title="Jaeger" **alt="Jaeger" width="60" height="60"/>&nbsp;
<img src="https://github.com/devicons/devicon/blob/master/icons/postgresql/postgresql-original.svg" title="PG" **alt="PG" width="60" height="60"/>&nbsp;
<img src="https://github.com/devicons/devicon/blob/master/icons/apachekafka/apachekafka-original.svg" title="Kafka" **alt="Kafka" width="60" height="60"/>&nbsp;
<img src="https://github.com/devicons/devicon/blob/master/icons/envoy/envoy-original.svg" title="Envoy" **alt="Envoy" width="60" height="60"/>&nbsp;
</p>

- **Go:** 1.24+
- **gRPC:** google.golang.org/grpc
- **OpenTelemetry (OTEL):**
  - **Metrics:** Prometheus → Grafana (визуализация метрик производительности)
  - **Logging:** uber-go/zap → Elasticsearch → Kibana (централизованное логирование и поиск)
  - **Tracing:** Jaeger (распределённая трассировка запросов)
- **Circuit Breaker:** github.com/sony/gobreaker
- **Database:** github.com/jackc/pgx (PostgreSQL driver)
- **Kafka:** github.com/IBM/sarama
- **Envoy:** github.com/envoyproxy/envoy
- **JWT:** github.com/golang-jwt/jwt

## Компоненты библиотеки

### Observability
#### Distributed Tracing (OpenTelemetry)
Полная интеграция с OpenTelemetry для распределённой трассировки запросов между микросервисами.
**Компоненты:**
- **gRPC Interceptor** - автоматическая инструментация gRPC вызовов
- **Metadata Carrier** - propagation trace context через gRPC metadata
- **Tracer** - настройка и инициализация OTEL tracer
**Экспорт:** OTLP (OpenTelemetry Protocol) в OTEL Collector
  
#### Metrics (Prometheus)
Сбор метрик производительности в формате Prometheus.
**Метрики:**
- **RPS** (Requests Per Second) - количество запросов в секунду
- **Latency Percentiles** - перцентили задержек (p50, p95, p99)
- **Error Rate** - процент ошибочных запросов
  
**Middleware:** автоматический сбор метрик через gRPC interceptor

#### Structured Logging (Zap)
Высокопроизводительное структурированное логирование на основе `uber-go/zap`.
**Особенности:**
- **Dual output** - одновременная запись в stdout и OTEL Collector
- **Custom OTEL Core** - кастомное ядро zap для отправки логов в OTEL
- **Контекстное логирование** - поддержка trace ID, request ID из context
- **Уровни логирования** - DEBUG, INFO, WARN, ERROR

### Отказоустойчивость
#### Circuit Breaker
Реализация паттерна Circuit Breaker на основе `gobreaker` с поддержкой half-open state.
**Возможности:**
- Автоматическое отключение недоступных сервисов при достижении порога ошибок
- Half-open state для проверки восстановления сервиса
- Настраиваемые параметры: `MaxRequests`, `Timeout`, `FailureRate`
- Логирование изменений состояния (Closed → Open → Half-Open)
**Конфигурация:**
- `ServiceName` - имя сервиса для идентификации
- `MaxRequest` - максимальное количество запросов в half-open state
- `Timeout` - время ожидания перед переходом в half-open
- `FailureRate` - порог ошибок для открытия circuit breaker (0.0-1.0)
#### Rate Limiter
Защита от перегрузки с ограничением количества запросов в секунду.
**Параметры:**
- `Limit` - максимальное количество запросов
- `Period` - временной интервал (обычно 1s)



### Middleware

#### gRPC Interceptors
**Circuit Breaker Middleware**
- Автоматическое применение circuit breaker к gRPC методам
- Возврат ошибки `UNAVAILABLE` при открытом circuit breaker

**Metrics Middleware**
- Сбор метрик для каждого gRPC метода (RPS, latency, errors)
- Автоматическое добавление labels (method, status_code)

**Rate Limiter Middleware**
- Ограничение частоты запросов на уровне gRPC interceptor
- Возврат ошибки `RESOURCE_EXHAUSTED` при превышении лимита

**Validation Middleware**
- Автоматическая валидация protobuf сообщений
- Конвертация ошибок валидации в gRPC статус коды

#### Kafka Middleware
**Logging Middleware**
- Логирование всех входящих/исходящих Kafka сообщений
- Трассировка обработки событий

### Clients
#### Database (PostgreSQL)
**Компоненты:**
- **PG Client** - обёртка над `pgx` для работы с PostgreSQL
- **Transaction Manager** - управление транзакциями с поддержкой context
- **Query Prettier** - форматирование SQL запросов для логирования
**Возможности:**
- Connection pooling
- Prepared statements
- Context-aware transactions
- Красивое логирование SQL запросов

#### Kafka
**Producer**
- Асинхронная отправка сообщений
- Автоматическая сериализация в protobuf
- Retry при ошибках отправки
**Consumer**
- Consumer Groups с балансировкой нагрузки
- Обработка сообщений с использованием handler pattern
- Автоматический commit offset

**Поддерживаемые события:**
- `user.created` - создание пользователя
- `user.deleted` - удаление пользователя

### Tokens & Authentication
#### JWT Generator
Генерация JWT токенов для аутентификации.
**Типы токенов:**
- **Access Token** - краткосрочный токен доступа
- **Refresh Token** - долгосрочный токен обновления

**Алгоритм:** HS256 (HMAC with SHA-256)

**Claims:**
- User ID
- User Email
- Expiration time (exp)
- Issued at (iat)

#### JWT Verifier
Проверка подписи и валидация JWT токенов.

**Функции:**
- Валидация подписи
- Проверка истечения срока действия
- Извлечение claims из токена

### Context Utilities

Набор утилит для работы с `context.Context`:

**Claims Context** - хранение JWT claims в контексте
```
// Сохранение claims
ctx = claimsctx.Set(ctx, claims)

// Извлечение claims
claims := claimsctx.Get(ctx)
```

**Trace ID Context** - trace ID для распределённой трассировки
```
ctx = traceIDctx.Set(ctx, traceID)
traceID := traceIDctx.Get(ctx)
```

**IP Context** - IP адрес клиента
```
ctx = ipctx.Set(ctx, clientIP)
ip := ipctx.Get(ctx)
```

**Transaction Context** - PostgreSQL транзакция
```
ctx = txctx.Set(ctx, tx)
tx := txctx.Get(ctx)
```

### Error Handling

#### Custom Error System
Кастомные ошибки с собственными статус кодами.

**Компоненты:**
- **Error Codes** - перечисление кодов ошибок (аналог HTTP status codes)
- **Error Type** - структура ошибки с кодом, сообщением и деталями
- **Converter** - маппинг внутренних кодов ошибок в gRPC статус коды

**Примеры кодов:**
- `USER_NOT_FOUND` → `codes.NotFound` → `grpc.NotFound`
- `INVALID_CREDENTIALS` → `codes.Unauthenticated` → `grpc.Unauthenticated`
- `PERMISSION_DENIED` → `codes.PermissionDenied` → `grpc.PermissionDenied`

#### Validation
Валидация protobuf сообщений с детальными ошибками.

**Возможности:**
- Автоматическая валидация через gRPC interceptor
- Поддержка protoc-gen-validate
- Человекочитаемые сообщения об ошибках

### Utilities

#### Closer
Graceful shutdown для корректного завершения работы сервисов.

**Использование:**
```
closer := closer.New()
closer.Add(func() error {
    return server.Shutdown()
})
closer.Add(func() error {
    return db.Close()
})

// При завершении работы
closer.Close()
```

## Infrastructure

Готовые docker-compose конфигурации для локальной разработки.

### Envoy Proxy

**Файлы:**
- `envoy.yaml` - конфигурация маршрутизации и load balancing
- `cert.pem` - TLS сертификат
- `messanger_descriptor.pb` - объединённый protobuf descriptor для gRPC-JSON transcoding

**Генерация descriptor:**
```
make generate-descriptor
```

**Возможности:**
- gRPC-JSON transcoding (HTTP → gRPC)
- Load balancing между инстансами сервисов
- TLS termination
- Rate limiting на уровне proxy

### Kafka

**Компоненты:**
- Kafka broker
- Kafka UI

**Топики:**
- `user.created`
- `user.deleted`

### OpenTelemetry Stack

**Компоненты:**
- **OTEL Collector** - сбор и экспорт метрик, трейсов и логов
- **Prometheus** - хранение метрик
- **Grafana** - визуализация метрик и трейсов

**Dashboards:**
- `dashboard.json` - общий dashboard для всех сервисов
- `dashboard_for_instanse.json` - dashboard для конкретного инстанса сервиса

**Endpoints:**
- Grafana: `http://localhost:3000`
- Prometheus: `http://localhost:9090`
- OTEL Collector: `http://localhost:4317` (gRPC), `http://localhost:4318` (HTTP)

## Версионирование

Библиотека использует **semantic versioning** (semver) с тэгами в GitHub.
**Формат:** `v1.3.0`
**Установка конкретной версии:**
```
go get github.com/WithSoull/platform_common@v1.3.0
```

## Использование
### Импорт модуля

```
import (
    "github.com/WithSoull/platform_common/pkg/logger"
    "github.com/WithSoull/platform_common/pkg/circuitbreaker"
    "github.com/WithSoull/platform_common/pkg/tokens/jwt"
)
```

### Пример: Logger с OTEL

```
import "github.com/WithSoull/platform_common/pkg/logger"

cfg := logger.Config{
    Level:       "INFO",
    AsJSON:      true,
    EnableOLTP:  true,
}

log := logger.New(cfg)
log.Info(ctx, "Service started", zap.String("service", "user-service"))
```

## Protocol Buffers

### События (proto/events/v1/events.proto)
Общие схемы событий для межсервисной коммуникации:

```
message UserCreated {
    string user_id = 1;
    string email = 2;
    int64 created_at = 3;
}

message UserDeleted {
    string user_id = 1;
    int64 deleted_at = 2;
}
```

## Архитектура зависимостей

```
┌─────────────────────────────────────────────┐
│           Microservices Layer               │
│   (UserServer, AuthService, ChatServer)     │
└─────────────────┬───────────────────────────┘
                  │ imports
┌─────────────────▼───────────────────────────┐
│      Messenger Platform Library             │
│                                             │
│  ┌──────────────┐  ┌──────────────────┐     │
│  │ Middleware   │  │  Observability   │     │
│  │ - Metrics    │  │  - Tracing       │     │
│  │ - Validation │  │  - Logging       │     │
│  │ - Auth       │  │  - Metrics       │     │
│  └──────────────┘  └──────────────────┘     │
│                                             │
│  ┌──────────────┐  ┌──────────────────┐     │
│  │ Resilience   │  │  Clients         │     │
│  │ - Circuit    │  │  - PostgreSQL    │     │
│  │   Breaker    │  │  - Kafka         │     │
│  │ - Rate       │  │  - JWT           │     │
│  │   Limiter    │  │                  │     │
│  └──────────────┘  └──────────────────┘     │
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │  Utilities (Context, Errors, Closer) │   │
│  └──────────────────────────────────────┘   │
└─────────────────┬───────────────────────────┘
                  │ uses
┌─────────────────▼───────────────────────────┐
│       Infrastructure (docker-compose)       │
│ Envoy | Kafka | OTEL | Prometheus | Grafana │
└─────────────────────────────────────────────┘
```

## Принципы дизайна

- **Separation of Concerns** - каждый компонент решает одну задачу
- **Dependency Injection** - все зависимости передаются через конструкторы
- **Interface-based** - работа через интерфейсы для тестируемости
- **Context-aware** - все операции поддерживают context.Context
- **Observability First** - встроенная поддержка метрик, трейсов и логов
