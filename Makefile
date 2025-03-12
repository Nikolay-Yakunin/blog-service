# TODO: Добавить нормальные переменные
build:
	go build -o main ./cmd

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

.PHONY: build test docker-up clean
