FROM golang:1.24.1-alpine AS deps

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN /go/bin/swag init -g cmd/blog-service/main.go -o docs --parseDependency --parseInternal --parseDepth 2
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/blog-service/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/main .
COPY templates ./templates
COPY static ./static
COPY --from=builder /app/docs ./docs
COPY config/config.yaml ./config.yaml

EXPOSE 8080
# Запускаем приложение
CMD ["./main"]
