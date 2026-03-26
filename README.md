# QuickSlot

Сервис онлайн-бронирования, объединяющий клиентов и бизнес в одном приложении. Система предоставляет удобное бронирование, отзывы и инструменты для роста компаний.

## Tech Stack

- Go 1.26
- MySQL 8.0
- JWT (HS256) для аутентификации
- golang-migrate для миграций
- Docker + Docker Compose

## Запуск

### С Docker

```bash
docker-compose up --build
```

Приложение будет доступно на `http://localhost:8080`.

### Локально

1. Создать MySQL базу данных `QuickSlot`
2. Скопировать `.env.example` в `.env` и заполнить данные
3. Запустить миграции вручную или через код
4. Запустить приложение:

```bash
go run ./cmd/app
```

## Переменные окружения

| Переменная | Описание | По умолчанию |
|---|---|---|
| DB_HOST | Хост MySQL | localhost |
| DB_PORT | Порт MySQL | 3306 |
| DB_USER | Пользователь MySQL | root |
| DB_PASS | Пароль MySQL | — |
| DB_NAME | Название базы | QuickSlot |
| DB_SSL | SSL-режим | false |
| JWT_SECRET | Секрет для JWT-токенов | — |
| SERVER_PORT | Порт сервера | 8080 |

## Структура проекта

```
QuickSlot/
├── cmd/app/             — точка входа
├── database/migrations/ — SQL-миграции
├── internal/
│   ├── handler/         — HTTP-хэндлеры
│   ├── middleware/       — Auth, RBAC, логирование, CORS
│   ├── model/           — структуры данных
│   ├── repository/      — работа с БД
│   ├── service/         — бизнес-логика
│   └── worker/          — фоновые задачи
├── pkg/
│   ├── auth/            — JWT утилиты
│   └── database/mysql/  — подключение к MySQL
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## API Endpoints

### Аутентификация

| Метод | Путь | Auth | Описание |
|---|---|---|---|
| POST | `/register` | — | Регистрация пользователя |
| POST | `/login` | — | Логин, возвращает JWT-токен |
| GET | `/me` | JWT | Проверка авторизации |

**Register:**
```json
POST /register
{
  "email": "user@example.com",
  "password": "pass1234"
}
```

**Login:**
```json
POST /login
{
  "email": "user@example.com",
  "password": "pass1234"
}
// Response: { "token": "eyJ..." }
```

### Организации

| Метод | Путь | Auth | Описание |
|---|---|---|---|
| POST | `/organizations/create` | JWT | Создать организацию |
| GET | `/organizations` | — | Список всех организаций |
| GET | `/organizations/get?id=1` | — | Получить организацию по ID |
| POST | `/organizations/update` | JWT + Admin | Обновить организацию |
| POST | `/organizations/delete` | JWT + Admin | Удалить организацию |

### Сотрудники

| Метод | Путь | Auth | Описание |
|---|---|---|---|
| POST | `/employees/create` | JWT + Admin | Добавить сотрудника |
| GET | `/employees?organization_id=1` | — | Сотрудники организации |
| POST | `/employees/update` | JWT + Admin | Обновить сотрудника |
| POST | `/employees/delete` | JWT + Admin | Удалить сотрудника |

### Слоты (расписание)

| Метод | Путь | Auth | Описание |
|---|---|---|---|
| POST | `/slots/generate` | JWT + Admin | Сгенерировать слоты для сотрудника |
| GET | `/slots/available?employee_id=1` | — | Доступные слоты |

**Generate slots:**
```json
POST /slots/generate
{
  "employee_id": 1,
  "date": "2026-04-01",
  "start_hour": 9,
  "end_hour": 18,
  "duration": 30
}
```

### Бронирование

| Метод | Путь | Auth | Описание |
|---|---|---|---|
| POST | `/appointments/book` | JWT | Забронировать слот |
| POST | `/appointments/cancel` | JWT | Отменить бронирование |
| GET | `/appointments/history` | JWT | История записей |

**Book:**
```json
POST /appointments/book
{ "slot_id": 1 }
```

**Cancel:**
```json
POST /appointments/cancel
{ "appointment_id": 1 }
```

**History с фильтрами:**
```
GET /appointments/history?from=2026-03-01&to=2026-04-01
```

### Отзывы

| Метод | Путь | Auth | Описание |
|---|---|---|---|
| POST | `/reviews/create` | JWT | Оставить отзыв |
| GET | `/reviews?organization_id=1` | — | Отзывы об организации |
| POST | `/reviews/delete` | JWT | Удалить свой отзыв |

**Create review:**
```json
POST /reviews/create
{
  "organization_id": 1,
  "rating": 5,
  "comment": "Отличный сервис!"
}
```

## Авторизация

Для защищённых эндпоинтов добавьте заголовок:
```
Authorization: Bearer <token>
```

Роли:
- `USER` — обычный пользователь (бронирование, отзывы)
- `ADMIN` — администратор (управление организациями, сотрудниками, слотами)

## Make-команды

| Команда | Описание |
|---|---|
| `make test` | Запуск всех тестов |
| `make test-v` | Тесты с подробным выводом (для демо) |
| `make test-cover` | Тесты + отчёт покрытия кода |
| `make build` | Собрать бинарь в `bin/` |
| `make run` | Запустить приложение локально |
| `make docker-up` | Поднять Docker |
| `make docker-down` | Остановить Docker |
| `make clean` | Удалить артефакты сборки |

## Тесты

```bash
# все тесты одной командой
make test-v

# тесты с покрытием
make test-cover
```

Покрытие: auth service, slot service, review service, JWT-пакет.

## Демонстрация через cURL

Ниже — полный сценарий работы с API через cURL.

### 1. Регистрация и логин

```bash
# Регистрация
curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@test.com","password":"pass1234"}'

# Логин — запомнить токен
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@test.com","password":"pass1234"}' | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "Token: $TOKEN"
```

### 2. Организации

```bash
# Создать организацию
curl -s -X POST http://localhost:8080/organizations/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Clinic Plus"}'

# Список организаций
curl -s http://localhost:8080/organizations
```

### 3. Слоты и бронирование

```bash
# Доступные слоты
curl -s "http://localhost:8080/slots/available?employee_id=1"

# Забронировать слот
curl -s -X POST http://localhost:8080/appointments/book \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"slot_id":1}'

# История бронирований
curl -s http://localhost:8080/appointments/history \
  -H "Authorization: Bearer $TOKEN"

# Отмена
curl -s -X POST http://localhost:8080/appointments/cancel \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"appointment_id":1}'
```

### 4. Отзывы

```bash
# Оставить отзыв
curl -s -X POST http://localhost:8080/reviews/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"organization_id":1,"rating":5,"comment":"Отличный сервис!"}'

# Отзывы организации
curl -s "http://localhost:8080/reviews?organization_id=1"
```

### 5. Проверка авторизации

```bash
# Без токена — получим 401
curl -s http://localhost:8080/me

# С токеном — получим 200
curl -s http://localhost:8080/me \
  -H "Authorization: Bearer $TOKEN"
```

## Фоновые задачи

- **Очистка просроченных слотов** — background goroutine с `time.Ticker` каждые 5 минут удаляет незабронированные слоты в прошлом. Используются channels для остановки при graceful shutdown.

## Graceful Shutdown

Сервер обрабатывает `SIGINT`/`SIGTERM` сигналы и корректно завершает:
1. Останавливает фоновый worker (через `close(done)`)
2. Завершает текущие HTTP-запросы (через `server.Shutdown`)
3. Таймаут завершения — 10 секунд

## Команда

- Amrayev Ruslan — Team Lead
- Darya Bashar — Core Backend Developer
- [Участник 3] — Scrum Master
- [Участник 4] — QA Engineer
