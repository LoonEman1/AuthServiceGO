# AuthServiceGO


Микросервис авторизации. Отвечает за юзеров, хеширование паролей и выдачу токенов.

Рефреш токены хранятся отдельной таблицей внутри бд, имеют дату создания и дату истечения. Используются миграции бд.



## Стек

* **Язык:** Go 1.25 (роутинг на стандартном `net/http`)

* **БД:** PostgreSQL (миграции + SQL запросы из библиотеки `sqlx`)

* **Токены:** JWT (библиотека `jwt/jwt/v5 v5.3.1`)

* **Деплой:** Docker + Docker Compose



## Ручки (API)

Все запросы работают с `Content-Type: application/json`.

| Метод | Эндпоинт | Описание | Body (что отправляем) | Ответ |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Регистрация нового юзера | `nickname`, `email`, `password`, `real_name` (может не указываться),  `birth_date` (может не указываться)| `201` + юзер, access и refresh токены |
| `POST` | `/api/v1/auth/login` | Вход в систему | `identifier` (Никнейм или почта), `password` | `200` + юзер, access и refresh токены |
| `POST` | `/api/v1/auth/refresh` | Обновление токенов | `refresh_token` | `200` + новые токены (access и refresh) |
| `POST` | `/api/v1/auth/logout` | Логаут | `refresh_token` | `204 No Content` (токен удален из БД) |



## Как запустить локально



1. Скопируй конфиг:

   ```bash
   cp .env.example .env
Заполни получившийся .env нужными тебе данными

2. В корне проекта выполни команду:

   ```bash
   docker-compose up --build
   ```

## Примеры тестовых запросов к серверу.

## Примеры (curl CMD)

Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d "{\"nickname\":\"testUser\",\"email\":\"test@example.com\",\"password\":\"testPass\",\"real_name\":\"Andrey\",\"birth_date\":\"2000-01-01T00:00:00Z\"}"
```

Login
  ```bash
curl -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d "{\"identifier\":\"testUser\",\"password\":\"testPass\"}"
```

Refresh
  ```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh -H "Content-Type: application/json" -d "{\"refresh_token\":\"ВСТАВЬ_ТОКЕН\"}"
```

Logout
  ```bash
curl -X POST http://localhost:8080/api/v1/auth/logout -H "Content-Type: application/json" -d "{\"refresh_token\":\"ВСТАВЬ_ТОКЕН\"}"
```
