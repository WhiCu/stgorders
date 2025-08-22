# Consumer Service

Сервис для обработки сообщений из Kafka с веб-интерфейсом и сохранением данных в PostgreSQL.

## Архитектура

Проект состоит из двух основных компонентов:
- **Kafka Consumer** - обрабатывает сообщения из Kafka топиков
- **Web Interface** - предоставляет HTTP API для работы с данными

## Структура проекта

```
consumer/
├── cmd/app/           # Точка входа в приложение
├── config/            # Конфигурационные файлы
├── db/               # База данных и миграции
├── internal/         # Внутренняя логика приложения
│   ├── app/         # Основное приложение
│   ├── config/      # Загрузка конфигурации
│   ├── kafka-consumer/  # Kafka потребитель
│   └── web-interface/   # Веб-интерфейс
├── pkg/              # Переиспользуемые пакеты
└── docker-compose.yml # Docker окружение
```

## Конфигурация (config.yaml)

### Server
```yaml
server:
  host: "localhost"        # Хост для веб-сервера
  port: "8080"            # Порт веб-сервера
  read_timeout: "10s"     # Таймаут чтения запроса
  write_timeout: "30s"    # Таймаут записи ответа
  idle_timeout: "30s"     # Таймаут простоя соединения
```

**Переменные окружения:**
- `SERVER_HOST` - хост сервера
- `SERVER_PORT` - порт сервера
- `SERVER_READ_TIMEOUT` - таймаут чтения
- `SERVER_WRITE_TIMEOUT` - таймаут записи
- `SERVER_IDLE_TIMEOUT` - таймаут простоя

### Logger
```yaml
logger:
  level: "info"           # Уровень логирования (debug, info, warn, error)
  path: "./logs/kafka-consumer.log"  # Путь к файлу логов
  size: 128               # Максимальный размер файла в МБ
  compress: false         # Сжатие старых логов
```

**Переменные окружения:**
- `LOGGER_LEVEL` - уровень логирования
- `LOGGER_PATH` - путь к файлу логов
- `LOGGER_SIZE` - максимальный размер файла в МБ
- `LOGGER_COMPRESS` - включить сжатие логов

**Особенности:**
- Если `LOGGER_PATH` пустой или не указан, логи выводятся только в stdout
- Поддерживаемые уровни: `debug`, `info`, `warn`, `error`

### Kafka
```yaml
kafka:
  brokers:                # Список Kafka брокеров
    - "localhost:9092"
  topic: "topic"     # Топик для чтения сообщений
  group_id: "group"    # ID группы потребителей
  worker_pool:            # Настройки пула воркеров
    size: 10              # Количество воркеров
    buf: 128              # Размер буфера сообщений
```

**Переменные окружения:**
- `KAFKA_BROKERS` - список брокеров (через запятую)
- `KAFKA_TOPIC` - название топика
- `KAFKA_GROUP_ID` - ID группы потребителей
- `KAFKA_WORKER_POOL_SIZE` - размер пула воркеров
- `KAFKA_WORKER_POOL_BUF` - размер буфера сообщений

**Особенности:**
- `KAFKA_BROKERS` принимает строку с брокерами через запятую (например: "localhost:9092,kafka:29092")

### Storage (PostgreSQL)
```yaml
storage:
  host: "localhost"       # Хост базы данных
  port: "5432"            # Порт базы данных
  user: "user"            # Пользователь БД
  password: "password"    # Пароль пользователя
  db: "db"          # Имя базы данных
```

**Переменные окружения:**
- `DB_HOST` - хост базы данных
- `DB_PORT` - порт базы данных
- `DB_USER` - пользователь базы данных
- `DB_PASSWORD` - пароль пользователя
- `DB_NAME` - имя базы данных

### Cache
```yaml
cache:
  size: 128               # Размер кэша в элементах
```

**Переменные окружения:**
- `CACHE_SIZE` - размер кэша в элементах

**Особенности:**
- Если `CACHE_SIZE = 0`, создается `NOPCache` (кэш без операций)
- Если `CACHE_SIZE > 0`, создается `LRUCache` с указанным размером
- Кэш автоматически вытесняет старые элементы при переполнении

## Переменные окружения

Все настройки можно переопределить через переменные окружения. Приоритет:
1. Переменные окружения
2. Файл конфигурации
3. Значения по умолчанию

**Основные переменные:**
```bash
# Путь к конфигурационному файлу
PATH_CONFIG=/path/to/config.yaml

# Настройки сервера
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Настройки логирования
LOG_LEVEL=debug
LOG_PATH=./logs/app.log

# Настройки базы данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=db

# Настройки Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=topic
KAFKA_GROUP_ID=group

# Настройки кэша
CACHE_SIZE=128
```

## Сборка и запуск

### Предварительные требования

- Go 1.24.3+
- Docker и Docker Compose
- Task (опционально, для упрощения команд)

### Запуск окружения

#### Docker Compose (рекомендуется)

Самый простой способ запуска всего окружения:

```bash
docker-compose up -d
```

**Особенности Docker Compose в проекте:**
- **PostgreSQL**: автоматическая инициализация через `init.sql`, health check
- **Kafka**: KRaft режим без Zookeeper, внутренний порт 29092 для контейнеров
- **Consumer**: многоэтапная сборка через Dockerfile, автоматические зависимости

**Особенности Dockerfile:**
- Многоэтапная сборка: `golang:1.25-alpine` → `distroless/static-debian12`
- Статическая компиляция с `CGO_ENABLED=0`
- Минимальный размер образа (без shell, пакетного менеджера)
- Non-root пользователь для безопасности

#### 🖥️ Локальный запуск

```bash
# Запуск только PostgreSQL
docker-compose up -d db

# Ожидание готовности БД
sleep 10

# Проверка статуса
docker-compose ps
```

### Сборка приложения

#### 🐳 Docker сборка

```bash
# Сборка образа
docker build -t consumer-service .

# Запуск контейнера
docker run -p 8080:8080 consumer-service
```

#### 🖥️ Локальная сборка

```bash
# Сборка
go build -o ./tmp/bin/main.exe ./cmd/app/main.go

# Или через Task
task build
```

### Запуск приложения

#### 🐳 Docker

```bash
# Запуск через docker-compose (все сервисы)
docker-compose up

# Запуск только приложения (требует запущенную БД и Kafka)
docker-compose up consumer
```

#### 🖥️ Локально

```bash
# Запуск собранного бинарника
./tmp/bin/main.exe

# Или запуск напрямую
go run ./cmd/app/main.go

# Или через Task
task run
```

### Тестирование

#### 🧪 GoConvey Framework

Проект использует **GoConvey** - мощный фреймворк для тестирования в стиле BDD (Behavior Driven Development).

#### Запуск тестов

```bash
# Стандартные Go тесты
go test ./... -v -cover

# Через Task
task test

# 🎯 GoConvey веб-интерфейс (рекомендуется для разработки)
goconvey

# GoConvey с параметрами
goconvey -host=0.0.0.0 -port=8081 -workDir=.
```

#### 🎯 GoConvey особенности

- **Веб-интерфейс**: автоматически открывает браузер с результатами тестов
- **BDD синтаксис**: `Convey("Given...", func() { ... })`
- **Автоматический перезапуск**: тесты перезапускаются при изменении кода
- **Покрытие кода**: показывает покрытие в реальном времени
- **Детальные отчеты**: полная информация о прохождении тестов

#### 📁 Структура тестов

```go
func TestHandler_ListenAndServe(t *testing.T) {
    Convey("Given a serviceable worker pool", t, func() {
        Convey("When context is cancelled immediately", func() {
            // ... setup code ...
            
            Convey("Then ListenAndServe should return nil without error", func() {
                err := h.ListenAndServe(ctx)
                So(err, ShouldBeNil)
            })
        })
    })
}
```

#### 🔍 Отдельные тесты

```bash
# Тесты конкретного пакета
go test ./internal/kafka-consumer/handler/... -v

# Тесты с покрытием
go test ./internal/kafka-consumer/handler/... -coverprofile=coverage.out

# Просмотр покрытия
go tool cover -html=coverage.out
```

## 🐳 Docker Compose

Файл `docker-compose.yml` содержит:

- **PostgreSQL 16** - основная база данных
  - Порт: 5432
  - Пользователь: user
  - Пароль: password
  - База: l0_test
  - Инициализация через `init.sql`
  - Health check для проверки готовности

- **Kafka** - брокер сообщений
  - Порт: 9092 (внешний), 29092 (внутренний)
  - KRaft режим (без Zookeeper)
  - Автоматическое создание топиков
  - Health check для проверки готовности

- **Consumer** - основное приложение
  - Собирается из Dockerfile
  - Порт: 8080
  - Зависит от готовности БД и Kafka
  - Использует переменные окружения для конфигурации

## 🐳 Dockerfile

Проект использует многоэтапную сборку для создания минимального образа:

### Этапы сборки

1. **Builder** (`golang:1.25-alpine`)
   - Установка зависимостей Go
   - Компиляция статического бинарного файла
   - Оптимизация размера через `-trimpath -ldflags "-s -w"`

2. **Runtime** (`gcr.io/distroless/static-debian12`)
   - Минимальный образ без shell и пакетного менеджера
   - Только необходимые системные библиотеки
   - Безопасность через non-root пользователя

### Особенности

- **Статическая компиляция**: `CGO_ENABLED=0` для переносимости
- **Минимальный размер**: использование distroless образа
- **Безопасность**: non-root пользователь по умолчанию
- **Оптимизация**: удаление отладочной информации из бинарного файла

### Переменные окружения для Docker

```bash
# Основные настройки
DB_HOST=db                    # Хост PostgreSQL
KAFKA_BROKERS=kafka:29092    # Внутренний адрес Kafka
LOG_LEVEL=debug              # Уровень логирования

# Путь к конфигурации
PATH_CONFIG=/src/config/config.yaml
```

### Запуск всех сервисов

```bash
# Сборка и запуск
docker-compose up --build

# Запуск в фоновом режиме
docker-compose up -d --build
```

### Запуск только базы данных

```bash
docker-compose up -d db
```

### Остановка

```bash
# Остановка всех сервисов
docker-compose down

# Остановка с удалением volumes
docker-compose down -v
```

### Просмотр логов

```bash
# Логи всех сервисов
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f consumer
docker-compose logs -f db
docker-compose logs -f kafka
```

## 📊 База данных

### Миграции

Миграции находятся в папке `db/migrations/`:
- `000001_create_users_table.up.sql` - создание таблицы пользователей
- `000001_create_users_table.down.sql` - откат изменений

### Схема

Основная схема БД описана в `db/schema.sql`

### SQLC

Проект использует SQLC для генерации Go кода из SQL запросов:
- Конфигурация: `db/sqlc.yaml`
- Запросы: `db/query.sql`
- Сгенерированный код: `db/pg/query.sql.go`

## 🔧 Разработка

### Hot Reload

Для разработки с автоматической перезагрузкой:

```bash
# Установка Air
go install github.com/cosmtrek/air@latest

# Запуск
air

# Или через Task
task air
```

### Логирование

Логи сохраняются в папку `logs/` с ротацией файлов.

### Переменные окружения

Поддерживается загрузка переменных из `.env` файла.

## 📝 API Endpoints

Веб-интерфейс предоставляет следующие эндпоинты:

- `GET /ping` - проверка работоспособности
- `GET /` - главная страница
- Другие эндпоинты для работы с заказами

## 🧪 Тестирование

Проект содержит тесты для всех основных компонентов:
- Unit тесты для сервисов
- Интеграционные тесты для обработчиков
- Тесты для кэша и хранилища

## 📋 Taskfile

Проект использует Task для автоматизации:

- `task test` - запуск тестов
- `task build` - сборка приложения
- `task run` - запуск приложения
- `task air` - запуск с hot reload
- `task up` - запуск миграций

## 🔍 Мониторинг

- Логирование в файлы с ротацией
- Метрики производительности
- Health check эндпоинты

## 🚨 Troubleshooting

### Проблемы с подключением к БД
1. Убедитесь, что PostgreSQL запущен: `docker-compose ps`
2. Проверьте настройки в `config.yaml`
3. Проверьте логи приложения

### Проблемы с Kafka
1. Убедитесь, что Kafka брокер доступен
2. Проверьте настройки топика и группы
3. Проверьте логи consumer'а

### Проблемы с портами
1. Убедитесь, что порты 8080 и 5432 свободны
2. Измените порты в `config.yaml` при необходимости
