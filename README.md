# Effective Mobile — Сервис управления подписками

REST API-сервис на Go и PostgreSQL для хранения и управления пользовательскими подписками.

## 🚀 Стек технологий

- Go (Golang)
- PostgreSQL
- Docker & Docker Compose
- REST API (JSON)
- [Chi Router](https://github.com/go-chi/chi)

---

## 📦 Запуск проекта

Убедитесь, что у вас установлен Docker и Docker Compose.

1. Клонируйте репозиторий:
   git clone https://github.com/dekatei/effective_mobile.git
   cd effective_mobile

2. Запустите контейнер
    docker-compose up --build
3. API будет доступен по адресу:
    http://localhost:8080

4. Примеры запросов  (или Postman)
    ➕ Добавить подписку

    curl -X POST http://localhost:8080/subscription \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "service_name": "vk",
    "price": 50,
    "start_date": "01-2022",
    "end_date": "01-2025"
  }'

  
    📋 Получить все подписки пользователя
    curl http://localhost:8080/subscription/60601fee-2bf1-4721-ae6f-7636e79a0cba


  
    📊 Расчёт стоимости подписок 'vk' за период
    curl "http://localhost:8080/cost/60601fee-2bf1-4721-ae6f-7636e79a0cba?start_date=07-2000&end_date=12-2030&service_name=vk"

    📊 Расчёт стоимости всех подписок за период
    curl "http://localhost:8080/cost/60601fee-2bf1-4721-ae6f-7636e79a0cba?start_date=07-2000&end_date=12-2030"
