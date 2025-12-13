 # Билд стадия
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости для компиляции
RUN apk add --no-cache git

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o idm-app ./cmd/api

# Финальная стадия
FROM alpine:latest

WORKDIR /root/

# Устанавливаем зависимости для работы
RUN apk --no-cache add ca-certificates tzdata

# Копируем бинарник из builder стадии
COPY --from=builder /app/idm-app .

# Порт приложения
EXPOSE 8080

# Запускаем приложение
CMD ["./idm-app"]