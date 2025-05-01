include .env

DOCKER_COMPOSE := docker-compose

#Запустить приложение(с Docker)
start:
	@echo "Запускаем контейнеры..."
	$(DOCKER_COMPOSE) up --build -d

# Остановить контейнеры
stop:
	@echo "Останавливаем контейнеры..."
	$(DOCKER_COMPOSE) down

# Запуск всех интеграционных тестов
test:
	@echo "Запускаем тесты..."
	go test ./tests/integration/... -bench=. -race

# Применить миграции (
migrate:
	@echo "Применяем миграции..."
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/ratelimiter?sslmode=disable" up

# Откатить миграции
migrate-down:
	@echo "Откатываем миграции..."
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/ratelimiter?sslmode=disable" down

# Запустить нагрузку Apache Bench с помощью Docker
ab-test:
	docker run jordi/ab -n 5000 -c 1000 http://host.docker.internal:8080/




