# g-store

## Архитектура 

- order service - логика заказов
- payment service - логика оплаты, счёта
- order notification service - оповещения об изменениях заказов
- api gateway - api шлюз для клиента
- frontend - клиентская часть

---

## Основные use cases

### Создание заказа

1. Заказ создается в `order service`
2. Публикуется событие `OrderCreated` в `order.events` - `outbox pattern` обеспечивает `at least once`
3. `payment service` - консьюмер `order.events`, `inbox pattern` + идемпотентная обработка + `outbox pattern` продьюсера обеспечивают `exactly once`  
4. `payment service` производит логику оплаты заказа, если все ок, деньги списываются со счёта 
5. `payment service` публикует событие о статусе оплаты `PaymentProcessed` в `payment.events` (`outbox pattern` обеспечивает `at least once`)
6. `order service` - консьюмер `payment.events`, обновляет статус заказа в зависимости от статуса оплаты
7. `order service` публикует событие `OrderStatusChanged` в `order.events` если оплата успешна, событие `OrderCancelled` иначе 

---

### Отмена заказа

1. Заказ отменяется если пользователь этого захотел, либо вынужден магазин, либо проблемы с оплатой (№7 в создании заказа)  
2. `order service` публикует событие `OrderCancelled` в `order.events` (`outbox, at least once`)
3. `payment service` консьюмит аналогично через `inbox-exactly-once`, если заказ уже был оплачен - деньги возвращаются на счёт
4. `order notification service` - косьюмер `order.events` - уведомляет клиентов-подписчиков об отмене

---

### Изменение статуса заказа

1. Изменение статуса заказа происходит со стороны персонала, либо при успешной оплате
2. `order service` публикует событие `OrderStatusChanged` в `order.events` (`outbox, at least once`)
3. `order notification service` - косьюмер `order.events` - уведомляет клиентов-подписчиков об изменении статуса

---

## Event-Driven Design

Данная архитектура была выбрана для последующего легкого внедрения новых подписчиков.

1. `order service` публикует события `OrderCreated`, `OrderCancelled`, `OrderStatusChanged`, консьюмит `PaymentProcessed`
2. `payment service` публикует событие `PaymentProcessed`, консьюмит `OrderCreated`, `OrderCancelled`
3. `order notification service` консьюмит `OrderCancelled`, `OrderStatusChanged`

---

## Взаимодействие сервисов

Взаимодействие асинхронное, производится через брокер `Apache Kafka`<br>
Для повышения производительности сообщения сериализуются в `protobuf`

---

## Stack
- `Go chi` - restful apis
- `Go gorrila` - websocket
- `Apache Kafka` - message broker
- `PostgreSQL` - databases
- `Nginx` - api gateway
- `Toastify js` - frontend
- `Docker` - containerisation

---

## Особенности реализации сервисов

1. Построены по Clean Architecture c элементами DDD
2. Есть базовая авторизация, разделение на роли staff / customer
3. Транзакции для поддержания консистентности данных
4. Логирование с фильтрами для анализа работоспособности

---

# Запуск

```sh
cd src
docker compose up --scale order-notification-service=3
```

1. Swagger для `order service` будет доступен на http://localhost/api/orders/swagger/index.html
2. Swagger для `payment service` будет доступен на http://localhost/api/payments/swagger/index.html
3. Frontend за 5 копеек (для `order-notification-service)` будет доступен на http://localhost:3000
4. Kafka-UI: http://localhost:9093

---