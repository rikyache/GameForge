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

Нужны Go 1.26, Docker Compose и CLI `golang-migrate`.

1. Клонируйте репозиторий и перейдите в каталог проекта:

   ```bash
   git clone https://github.com/rikyache/GameForge.git
   cd GameForge
   ```

2. Создайте локальный `.env` по образцу [.env.example](.env.example) и замените значения-заглушки. Не добавляйте `.env` в Git. Для запуска Go на компьютере укажите `DB_HOST=localhost` и `REDIS_HOST=localhost`. Docker Compose передаёт контейнеру приложения адреса сервисов самостоятельно.

3. Поднимите PostgreSQL и Redis:

   ```bash
   docker compose up -d postgres redis
   ```

4. Примените миграции к базе, указанной в `DB_NAME`. Вместо `<DATABASE_URL>` подставьте локальный PostgreSQL URL с параметрами из `.env`:

   ```bash
   migrate -path ./migrations -database "<DATABASE_URL>" up
   ```

5. Запустите приложение одним из способов:

   ```bash
   go run .
   ```

   или

   ```bash
   docker compose up --build -d app
   ```

API доступно по адресу `http://localhost:8080`.

## Тесты

Интеграционные тесты в `internal/repository` используют отдельную базу `practice_test`. Скрипт Docker создаёт её при первой инициализации тома PostgreSQL. Для уже существующего тома проверьте, что база создана. Тесты берут `DB_HOST`, `DB_PORT`, `DB_USER` и `DB_PASSWORD` из окружения или корневого `.env`; имя тестовой базы зафиксировано в тестовом коде.

После запуска PostgreSQL примените те же миграции к тестовой базе и запустите тесты:

```bash
migrate -path ./migrations -database "<TEST_DATABASE_URL>" up
go test ./...
```

`<TEST_DATABASE_URL>` — локальный PostgreSQL URL для `practice_test` с теми же пользователем и паролем, что в `.env`. Тесты очищают таблицы `player_games`, `players` и `games` и сбрасывают их идентификаторы. Не храните в `practice_test` нужные данные.

Сейчас интеграционные тесты проверяют часть операций репозиториев игроков, игр и покупок. Авторизация, возврат, кэш и HTTP-обработчики пока не покрыты тестами.

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
