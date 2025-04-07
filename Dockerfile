FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/blog-service
# Проверяем наличие .env файла и копируем .env.example в .env, если файл .env не существует
RUN if [ ! -f .env ]; then cp .env.example .env || touch .env; fi

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]
