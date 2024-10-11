# Указываем базовый образ
FROM golang:1.23-alpine AS build

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum, чтобы установить зависимости
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы в контейнер
COPY . .

# Собираем приложение
RUN go build -o SongsLibrary ./cmd/main.go

# Используем минимальный образ для запуска приложения
FROM alpine:latest

# Устанавливаем bash
RUN apk add --no-cache bash

# Устанавливаем рабочую директорию для запуска
WORKDIR /root/

# Копируем бинарник приложения из предыдущего образа
COPY --from=build /app/SongsLibrary .

# Копируем wait-for-it скрипт
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Открываем порт, который будет использован
EXPOSE ${APP_PORT}

# Указываем команду для запуска приложения с ожиданием базы данных
CMD ["bash", "/wait-for-it.sh", "db:5432", "--", "./SongsLibrary"]