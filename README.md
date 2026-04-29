# AuthServiceGO


Микросервис авторизации. Отвечает за юзеров, хеширование паролей и выдачу токенов.

Рефреш токены хранятся отдельной таблицей внутри бд, имеют дату создания и дату истечения. Используются миграции бд. Используется кафка для отправки в топик сообщения о регистрации пользователя.

Данный сервис используется как Git submodule в [project-deploy](https://github.com/LoonEman1/project-deploy) для запуска mailer-сервиса и сервиса авторизации в едином окружении.

## Стек

* **Язык:** Go 1.25 (роутинг на стандартном `net/http`)

* **БД:** PostgreSQL (миграции + SQL запросы из библиотеки `sqlx`)

* **Токены:** JWT (библиотека `jwt/jwt/v5 v5.3.1`)

* **Деплой:** Docker + Docker Compose



## Ручки (API)

Все запросы работают с `Content-Type: application/json`.

| Метод | Эндпоинт | Описание | Body (что отправляем) | Ответ |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Регистрация юзера | `nickname`, `email`, `password`, `real_name`*, `birth_date`* | `201` + Почта юзера на которую отправится код (код отправляется в кафку, для обработки другим микросервисом, живет 15 минут)|
| `POST` | `/api/v1/auth/verify` | Подтверждение почты | `email`, `code` (код, который пришел на почту юзеру) | `200` + Юзер + Токены |
| `POST` | `/api/v1/auth/verify/resend` | Отправить повторно код на почту | `email` | `204 No Content` |
| `POST` | `/api/v1/auth/login` | Вход в систему | `identifier` (Nick/Email), `password` | `200` + Юзер + Токены |
| `POST` | `/api/v1/auth/refresh` | Обновление токенов | `refresh_token` | `200` + Новые токены |
| `POST` | `/api/v1/auth/logout` | Выход из системы | `refresh_token` | `204 No Content` |



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
Для подтверждения необходимо получить код, сгенерированный при регистрации или при /verify/resend
Просмотреть локально сгенерированный код можно при помощи команды
  ```bash
  docker exec -it authservice-kafka-1 /opt/kafka/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic email-notifications --from-beginning
  ```

Verify
  ```bash
  curl -X POST http://localhost:8080/api/v1/auth/verify -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"code\":\"код из кафки\"}"
```

Resend

  ```bash
curl -X POST http://localhost:8080/api/v1/auth/verify/resend -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\"}"
```
