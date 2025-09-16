// Package config предоставляет структуры и функции для загрузки и управления
// конфигурацией приложения из файлов (например, YAML) и переменных окружения.
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config содержит все конфигурационные параметры приложения.
type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"`
	Database   DBConfig
	HTTPServer HTTPServerConfig
	Auth       AuthConfig
	Minio      MinioConfig
	SMS        SMSConfig
}

// DBConfig содержит параметры для подключения к базе данных.
type DBConfig struct {
	URL string `yaml:"url" env:"DATABASE_URL" env-required:"true"`
}

// HTTPServerConfig содержит параметры для HTTP-сервера.
type HTTPServerConfig struct {
	Port        string        `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

// AuthConfig содержит параметры для аутентификации (JWT).
type AuthConfig struct {
	JWTSecretKey string        `yaml:"jwt_secret_key" env:"JWT_SECRET_KEY" env-required:"true"`
	TokenTTL     time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"1h"`
}

// MinioConfig содержит параметры для подключения к S3-совместимому хранилищу MinIO.
type MinioConfig struct {
	Endpoint   string `yaml:"endpoint" env:"MINIO_ENDPOINT" env-required:"true"`
	AccessKey  string `yaml:"access_key" env:"MINIO_ACCESS_KEY" env-required:"true"`
	SecretKey  string `yaml:"secret_key" env:"MINIO_SECRET_KEY" env-required:"true"`
	BucketName string `yaml:"bucket_name" env:"MINIO_BUCKET_NAME" env-required:"true"`
	UseSSL     bool   `yaml:"use_ssl" env:"MINIO_USE_SSL" env-default:"false"`
}

// SMSConfig содержит параметры для интеграции с SMS-шлюзом.
type SMSConfig struct {
	APIKey     string `yaml:"api_key" env:"SMS_API_KEY" env-required:"true"`
	SenderName string `yaml:"sender_name" env:"SMS_SENDER_NAME" env-required:"true"`
}

// MustLoad загружает конфигурацию.
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "internal/config/config.yml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file not found at path: %s. Relying on environment variables.", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
