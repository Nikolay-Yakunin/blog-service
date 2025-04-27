# TODO: Добавить нормальные переменные
# Переменные для миграций
MIGRATE_CMD = migrate
DATABASE_URL ?= $(shell grep DATABASE_URL .env | cut -d '=' -f2-)
MIGRATIONS_PATH = migrations

build:
	go build -o main ./cmd/blog-service/

# Установка migrate CLI (если не установлен)
install-migrate:
	@command -v migrate >/dev/null 2>&1 || \
		(echo "Installing golang-migrate/migrate..."; \
		 go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)

# Применить все доступные миграции
migrate-up: install-migrate
	@echo "Applying migrations from $(MIGRATIONS_PATH) to $(DATABASE_URL)"
	$(MIGRATE_CMD) -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

# Откатить последнюю примененную миграцию
migrate-down: install-migrate
	@echo "Rolling back last migration from $(MIGRATIONS_PATH) on $(DATABASE_URL)"
	$(MIGRATE_CMD) -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

# Откатить все миграции
migrate-down-all: install-migrate
	@echo "Rolling back all migrations from $(MIGRATIONS_PATH) on $(DATABASE_URL)"
	$(MIGRATE_CMD) -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down

# Создать новую пару файлов миграции (например, make migrate-create name=add_feature)
migrate-create:
	@read -p "Enter migration name: " name; \
	 $(MIGRATE_CMD) create -ext sql -dir migrations -seq $$name

# Показать текущий статус миграций
migrate-status: install-migrate
	@echo "Checking migration status for $(DATABASE_URL)"
	$(MIGRATE_CMD) -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version

# Запуск тестов с генерацией отчёта покрытия
test:
	go test -coverprofile=coverage.out ./...

# Запуск docker-compose (с пересборкой образов)
docker-up:
	docker-compose up --build
# TODO: Добавить стадию для lcov

# Очистка сгенерированных файлов
clean:
	rm -f main coverage.out

.PHONY: build test docker-up clean swagger install-migrate migrate-up migrate-down migrate-down-all migrate-create migrate-status

swagger:
	swag init -g cmd/blog-service/main.go -o docs --parseDependency --parseInternal --parseDepth 2
	go run cmd/app/main.go
