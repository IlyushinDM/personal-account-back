# --- Этап загрузки migrate ---
FROM golang:1.24-alpine AS migrate_builder
RUN apk add --no-cache git
# Устанавливаем migrate с поддержкой postgres
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


# --- Этап сборки ---
FROM golang:1.24-alpine AS builder

# Копируем бинарный файл migrate из предыдущего этапа
COPY --from=migrate_builder /go/bin/migrate /usr/local/bin/migrate

# Устанавливаем пакет с данными о часовых поясах
RUN apk add --no-cache tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server -ldflags "-w -s" ./cmd/app

# --- Финальный этап ---
FROM alpine:latest

COPY --from=builder /app/server /server
# Копируем утилиту migrate в финальный образ, чтобы сервис migrate мог её использовать
COPY --from=builder /usr/local/bin/migrate /migrate
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV ZONEINFO /usr/share/zoneinfo

# Копируем миграции
COPY ./migrations /migrations

EXPOSE 8080
# Эта точка входа будет использоваться только сервисом 'app'
ENTRYPOINT ["/server"]