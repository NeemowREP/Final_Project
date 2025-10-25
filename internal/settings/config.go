package settings

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port         string        `envconfig:"TODO_PORT" default:"7540"`
	LogLevel     string        `envconfig:"LOG_LEVEL" default:"info"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
	WebDir       string        `envconfig:"WEB_DIR" default:"../../web"`
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env файл не найден, используются default настройки")
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Ошибка загрузки настроек: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка получения текущей директории")
	}
	cfg.WebDir = filepath.Join(cwd, cfg.WebDir)
	log.Printf("Используется web-директория: %s", cfg.WebDir)
	return cfg
}
