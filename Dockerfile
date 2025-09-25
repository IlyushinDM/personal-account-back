# --- Этап сборки ---
FROM golang:1.24-alpine AS builder

# Устанавливаем пакет с данными о часовых поясах
RUN apk add --no-cache tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server -ldflags "-w -s" ./cmd/app

# --- Финальный этап ---
FROM scratch

# Копируем бинарник из этапа сборки
COPY --from=builder /app/server /server
# Копируем базу данных часовых поясов
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Устанавливаем переменную окружения, чтобы Go знал, где искать tzdata
ENV ZONEINFO /usr/share/zoneinfo

# Копируем миграции
COPY ./migrations /migrations

EXPOSE 8080
ENTRYPOINT ["/server"]