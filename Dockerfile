# Базовый образ Go
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем приложение
RUN go build -o nail_bot cmd/bot/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник и файлы
COPY --from=builder /app/nail_bot .
COPY --from=builder /app/fonts ./fonts
COPY --from=builder /app/conf ./conf

# Устанавливаем зависимости для работы
RUN apk add --no-cache ca-certificates tzdata

# Переменные окружения (переопределяются при запуске)
ENV TELEGRAM_BOT_TOKEN=""
ENV DB_CONNECTION_STRING=""
ENV ADMIN_ID=""

# Запускаем бота
ENTRYPOINT ["./nail_bot"]
