FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download

# Устанавливаем swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Генерируем документацию Swagger (путь к main.go исправлен)
# Используем /go/bin/swag, так как go install помещает бинарник туда
RUN /go/bin/swag init -g cmd/app/main.go -o docs --parseDependency --parseInternal --parseDepth 2

# Исправлено: компилируем cmd/app/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go
# Проверяем наличие .env файла и копируем .env.example в .env, если файл .env не существует
# RUN if [ ! -f .env ]; then cp .env.example .env || touch .env; fi
# Убрал копирование .env из builder, т.к. он должен монтироваться или передаваться при запуске

FROM alpine:latest
RUN apk --no-cache add ca-certificates
# Определяем рабочую директорию
WORKDIR /app
# Копируем бинарный файл
COPY --from=builder /app/main .
# Копируем шаблоны и статику
COPY templates ./templates
COPY static ./static
# Копируем сгенерированную документацию
COPY --from=builder /app/docs ./docs

# Убрал COPY .env - конфигурацию лучше передавать через docker-compose environment или volumes

EXPOSE 8080
# Запускаем приложение
CMD ["./main"]
