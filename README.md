# Shelf

**Shelf** — backend-сервис для управления каталогом игр и пользовательской библиотекой.

Проект реализует регистрацию и авторизацию пользователей, работу с балансом, добавление игр в библиотеку, удаление и возврат игр, а также кэширование данных через Redis.

## Возможности

- регистрация и авторизация пользователей;
- JWT-аутентификация;
- logout с отзывом токена через Redis;
- просмотр профиля пользователя;
- пополнение баланса;
- просмотр каталога игр;
- добавление игр в библиотеку;
- удаление игр из библиотеки;
- возврат игр;
- кэширование данных пользователей и библиотек;
- логирование HTTP-запросов.

## Стек

- Go
- `net/http`
- PostgreSQL
- Redis
- JWT
- Docker / Docker Compose
- golang-migrate

## Архитектура

Приложение разделено на несколько слоёв:

```text
Handler
   ↓
Service
   ↓
Repository
   ↓
PostgreSQL
```

Для части запросов используется кэширующий слой поверх репозиториев:

```text
Service
   ↓
Cached Repository
   ↓
Repository
   ↓
PostgreSQL
        ↕
      Redis
```

HTTP handlers отвечают за работу с запросами и ответами.

Service layer содержит бизнес-логику приложения.

Repository layer отвечает за работу с PostgreSQL.

Redis используется для кэширования данных и хранения отозванных JWT-токенов.

## Аутентификация

После успешного входа пользователь получает JWT-токен.

Для доступа к защищённым endpoint'ам токен необходимо передавать в заголовке:

```http
Authorization: Bearer <token>
```

При logout идентификатор токена (`jti`) сохраняется в Redis до момента истечения его срока действия.

Middleware при каждом защищённом запросе проверяет:

- корректность JWT;
- срок действия токена;
- наличие токена в blacklist.

## API

### Авторизация

```http
POST /auth/register
POST /auth/login
POST /auth/logout
```

### Пользователь

```http
GET  /me
POST /me/deposit
```

### Библиотека пользователя

```http
GET    /me/games
POST   /me/games
DELETE /me/games/{gameID}
POST   /me/games/{gameID}/refund
```

### Игры

```http
GET    /games
GET    /games/{id}
POST   /games
DELETE /games/{id}
```

### Пользователи

```http
GET    /players
GET    /players/{id}
POST   /players
DELETE /players/{id}
```

## Структура проекта

```text
.
├── internal
│   ├── auth
│   ├── cache
│   ├── database
│   ├── handlers
│   ├── middleware
│   ├── repository
│   └── service
│
├── migrations
├── docker-compose.yml
├── go.mod
└── main.go
```

### `auth`

Работа с JWT и blacklist токенов.

### `cache`

Redis-клиент и реализация кэша.

### `database`

Подключение к PostgreSQL.

### `handlers`

HTTP handlers приложения.

### `middleware`

Middleware для аутентификации и логирования.

### `repository`

Работа с хранилищем данных.

### `service`

Бизнес-логика приложения.

## Запуск

### 1. Клонировать репозиторий

```bash
git clone <repository-url>
cd shelf
```

### 2. Настроить переменные окружения

Создать `.env` файл и указать настройки PostgreSQL, Redis и JWT.

Пример:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=shelf

REDIS_HOST=localhost
REDIS_PORT=6379

JWT_SECRET=secret
```

### 3. Запустить инфраструктуру

```bash
docker compose up -d
```

### 4. Запустить приложение

```bash
go run .
```

После запуска API будет доступно по адресу:

```text
http://localhost:8080
```

## Что реализовано в проекте

В рамках проекта я практиковал:

- построение REST API на стандартном `net/http`;
- разделение приложения на handler/service/repository;
- работу с PostgreSQL;
- использование Redis для кэширования;
- JWT-аутентификацию;
- middleware;
- работу с `context.Context`;
- отзыв JWT-токенов через blacklist;
- транзакционную бизнес-логику;
- работу с Docker.

## Планы

- добавить OpenAPI / Swagger;
- увеличить покрытие тестами;
- улучшить обработку ошибок;
- добавить роли пользователей;
- переработать управление каталогом игр;
- добавить более сложные сценарии работы с покупками и возвратами.