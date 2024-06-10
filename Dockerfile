# Используем официальный образ Go для сборки приложения
FROM golang:latest AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем все файлы в рабочую директорию
COPY . .

# Скачиваем зависимости
RUN go mod tidy
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp

# Используем минимальный образ для запуска приложения
FROM alpine:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /root/

# Устанавливаем зависимости (если нужны)
RUN apk --no-cache add ca-certificates

# Копируем бинарный файл и конфигурационный файл из предыдущего этапа
COPY --from=builder /app/myapp .
COPY --from=builder /app/config/config.json ./config/config.json

# Устанавливаем права на выполнение файла
RUN chmod +x myapp

# Открываем порт, на котором будет работать приложение
EXPOSE 8080

# Команда для запуска приложения
CMD ["./myapp"]
