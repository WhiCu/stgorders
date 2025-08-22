# Consumer Service

[![Go Version](https://img.shields.io/badge/Go-1.24.3+-blue.svg)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/WhiCu/stgorders)](https://goreportcard.com/report/github.com/WhiCu/stgorders)

Сервис для обработки сообщений из Kafka с веб-интерфейсом и сохранением данных в PostgreSQL и кэше.

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
- Многоэтапная сборка: `golang:1.25-alpine` → `alpine:3.22.0`
- Статическая компиляция с `CGO_ENABLED=0`
- Минимальный размер образа (без shell, пакетного менеджера)

#### Локальная сборка

```bash
# Сборка
go build -o ./tmp/bin/main.exe ./cmd/app/main.go

# Или через Task
task build
```

### Запуск приложения

#### Локально

```bash
# Запуск собранного бинарника
./tmp/bin/main.exe

# Или запуск напрямую
go run ./cmd/app/main.go

# Или через Task
task run
```

### Тестирование

#### GoConvey 

Проект использует **GoConvey**

#### Запуск тестов

```bash
# Стандартные Go тесты
go test ./... -v -cover

# Через Task
task test

# GoConvey веб-интерфейс
goconvey
```

### Логирование

Логи сохраняются в папку `logs/` с ротацией файлов (LOGGER_PATH != "").

### Переменные окружения

Поддерживается загрузка переменных из `.env` файла.

## API Endpoints

Веб-интерфейс предоставляет следующие эндпоинты:

- `GET /ping` - проверка работоспособности
- `GET /` - главная страница
- `GET /:orderUID` - получение заказа по UID