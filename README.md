# Consumer Service

[![Go Version](https://img.shields.io/badge/Go-1.24.3+-blue.svg)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/WhiCu/stgorders)](https://goreportcard.com/report/github.com/WhiCu/stgorders)

Сервис для обработки сообщений из Kafka с веб-интерфейсом и сохранением данных в PostgreSQL и кэше.

## Интерфейс

![kafka-consumer interface](/.github/kafka-consumer.gif)

## Архитектура

Проект состоит из двух основных компонентов:
- **Kafka Consumer** - обрабатывает сообщения из Kafka топиков
- **Web Interface** - предоставляет HTTP API для работы с данными

## Структура проекта

```
consumer/
├── cmd/                # Точки входа в приложение
│   ├── app/           # Основное приложение (веб-сервер + Kafka consumer)
│   ├── cli/           # CLI интерфейс с подкомандами (cobraCLI)
│   ├── kafka/         # Утилиты для работы с Kafka
│   └── migrate/       # Утилиты для миграций БД
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

### Основная команда

```bash
# Запуск основного приложения (веб-сервер + Kafka consumer)
stgorders

# С флагом для создания Kafka топика
stgorders --topic
# или
stgorders -t
```

### Подкоманды

#### `stgorders app`
Запуск основного приложения (веб-сервер + Kafka consumer):
```bash
stgorders app
```

#### `stgorders migrate`
Запуск миграций базы данных:
```bash
stgorders migrate
```

#### `stgorders topic`
Создание Kafka топика:
```bash
stgorders topic
```

### Сборка CLI

```bash
# Сборка CLI приложения
go build -o ./tmp/bin/stgorders.exe ./cmd/cli

# Или через Task
task build        # Сборка основного приложения
task build:cli    # Сборка CLI приложения
```

## Конфигурация (config.yaml)

### Server
```yaml
server:
  host: "localhost"       # Хост для веб-сервера
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
  level: "info"                       # Уровень логирования (debug, info, warn, error)
  path: "./logs/kafka-consumer.log"   # Путь к файлу логов
  size: 128                           # Максимальный размер файла в МБ
  compress: false                     # Сжатие старых логов
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
  topic: "topic"          # Топик для чтения сообщений
  group_id: "group"       # ID группы потребителей
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
  db: "db"                # Имя базы данных
```

**Переменные окружения:**
- `DB_HOST` - хост базы данных
- `DB_PORT` - порт базы данных
- `DB_USER` - пользователь базы данных
- `DB_PASSWORD` - пароль пользователя
- `DB_NAME` - имя базы данных

### Migate
```yaml
migrate:
  dir: "./db/migrations"         # Директория с файлами миграции     
```

**Переменные окружения:**
- `MIGRATE_DIR` - директория с файлами миграции

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
# Путь к конфигурационному файлу (по умолчанию: config/config.yaml)
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

**Примечание:** При использовании CLI (`stgorders`) миграции БД запускаются автоматически при старте приложения.

**Особенности Dockerfile:**
- Многоэтапная сборка: `golang:1.25-alpine` → `alpine:3.22.0`
- Статическая компиляция с `CGO_ENABLED=0`
- Минимальный размер образа (без shell, пакетного менеджера)
- Сборка CLI приложения с cobraCLI

### Запуск приложения

#### Локально

```bash

# Или запуск напрямую
go run ./cmd/app/main.go

# Или через Task
task run          # Запуск основного приложения
task run:cli      # Запуск CLI приложения

# Запуск через CLI (рекомендуется)
./tmp/bin/stgorders.exe

# Или запуск CLI напрямую
go run ./cmd/main.go
```

### Тестирование

#### GoConvey 

Проект использует [**GoConvey**](https://github.com/smartystreets/goconvey)

#### Запуск тестов

```bash
# Стандартные Go тесты
go test ./... -v -cover

# Через Task
task test

# GoConvey веб-интерфейс
goconvey
```

### Benchmark

#### bombardier  

Проект использует [**bombardier**](https://github.com/codesenberg/bombardier)

#### Cache ON

```bash
Bombarding http://localhost:8080/order/b563feb7b2b84b6test1 with 1000000 request(s) using 125 connection(s)
 1000000 / 1000000 [=================================================] 100.00% 23407/s 42s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec     23428.77    5937.70   35345.12
  Latency        5.33ms    10.90ms      1.13s
  HTTP codes:
    1xx - 0, 2xx - 1000000, 3xx - 0, 4xx - 0, 5xx - 0
    others - 0
  Throughput:    23.39MB/s
```

#### Cache OFF

```bash
Bombarding http://localhost:8080/order/b563feb7b2b84b6test1 with 1000000 request(s) using 125 connection(s)
 1000000 / 1000000 [=================================================] 100.00% 11850/s 1m24s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec     11879.31    2345.07   33052.88
  Latency       10.53ms     5.93ms   720.88ms
  HTTP codes:
    1xx - 0, 2xx - 1000000, 3xx - 0, 4xx - 0, 5xx - 0
    others - 0
  Throughput:    11.85MB/s
```

### Логирование

Логи сохраняются в папку `logs/kafka-consumer.log` с ротацией файлов (*Если LOGGER_PATH  задан по умолчанию*).


## API Endpoints

Веб-интерфейс предоставляет следующие эндпоинты:

- `GET /ping` - проверка работоспособности
- `GET /` - главная страница
- `GET /:orderUID` - получение заказа по UID в формате json