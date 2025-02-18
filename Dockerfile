# Stage 1: Сборка приложения
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем исполняемый файл
RUN CGO_ENABLED=0 GOOS=linux go build -o SongLibraryApi main.go

# Stage 2: Финальный образ
FROM alpine:latest

WORKDIR /app

# Устанавливаем необходимые утилиты
RUN apk add --no-cache bash

# Создаём директорию для логов
RUN mkdir -p log

# Копируем исполняемый файл из первого этапа
COPY --from=builder /app/SongLibraryApi .

# Копируем миграции и конфигурацию
COPY migrations ./migrations
COPY .env .

# Запускаем приложение
CMD ["./SongLibraryApi"]