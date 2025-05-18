# Основной функционал

* Регистрация и логин пользователей (JWT)
* Создание и фильтрация ПВЗ
* Приёмка товаров: создание, закрытие, добавление/удаление товаров
* gRPC метод получения всех ПВЗ без авторизации
* Метрики Prometheus на порту 9000
* OpenAPI спецификация REST API
* CI/CD pipeline (GitHub Actions)

---

## Структура проекта

```plaintext
.
├── cmd
│   └── pvz_application
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   ├── pvz_server.go
│   │   ├── setup_routers.go
│   │   └── transport
│   │       ├── grpc
│   │       │   ├── pvz.pb.go
│   │       │   ├── pvz_grpc.pb.go
│   │       │   └── pvz_handler.go
│   │       └── rest
│   │           ├── middleware
│   │           │   ├── auth_middleware.go
│   │           │   └── metrics.go
│   │           ├── product_controller.go
│   │           ├── product_dto.go
│   │           ├── prometheus_metrics
│   │           │   └── prometheus.go
│   │           ├── pvz_controller.go
│   │           ├── pvz_dto.go
│   │           ├── reception_controller.go
│   │           ├── reception_dto.go
│   │           ├── user_controller.go
│   │           └── user_dto.go
│   ├── db
│   │   ├── PostgreSQL
│   │   │   ├── product_repositoryImpl.go
│   │   │   ├── pvz_repositoryImpl.go
│   │   │   ├── reception_repositoryImpl.go
│   │   │   └── user_repositoryImpl.go
│   │   ├── db.go
│   │   ├── product_repository.go
│   │   ├── pvz_repository.go
│   │   ├── reception_repository.go
│   │   ├── sq
│   │   │   └── sqlbuilder.go
│   │   └── user_repository.go
│   ├── domain
│   │   └── models
│   │       ├── product_model.go
│   │       ├── pvz_model.go
│   │       ├── reception_model.go
│   │       └── user_model.go
│   └── logger
│       └── logger.go
└── tests
    ├── integration
    │   └── controller
    │       └── user_controller_test.go
    └── perfomarnce
        └── load_test.go
```

---

## Тестирование

```bash
go test tests/integration/controller/ -v
```

---

## gRPC

* Интерфейс описан в `internal/app/transport/grpc/*.proto`
* Сервер доступен на порту **3000**
* Авторизация не требуется

---

## Метрики

* Настроены сервисом Prometheus
* Доступны по эндпоинту **`/metrics`**
* Порт **9000**

---

## Docker (todo)

* В файлах Dockerfile и docker-compose.yml
* Запуск:
```bash
docker-compose up --build
```

---

## Схема базы данных (todo)

В разработке
