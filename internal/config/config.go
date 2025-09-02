package config

import (
	"log"
	"os"

	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer HTTPServer `yaml:"http_server"`
	DB         DB         `yaml:"db"`
	Kafka      Kafka      `yaml:"kafka"`
	Cache      Cache      `yaml:"cache"`
	Log        Log        `yaml:"log"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"HTTP_ADDRESS" env-default:":8081"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type DB struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5433"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"dbname" env:"DB_NAME" env-default:"L0_db"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
}

type Kafka struct {
	Brokers  []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:","`
	Topic    string   `yaml:"topic" env:"KAFKA_TOPIC" env-default:"orders"`
	GroupID  string   `yaml:"group_id" env:"KAFKA_GROUP_ID" env-default:"orderservice"`
	DLQTopic string   `yaml:"dlq_topic" env:"KAFKA_DLQ_TOPIC" env-default:"orders-dlq"`
	ClientID string   `yaml:"client_id" env:"KAFKA_CLIENT_ID" env-default:"orderservice"`
}

type Cache struct {
	MaxEntries int `yaml:"max_entries" env:"CACHE_MAX_ENTRIES" env-default:"1000"`
}

type Log struct {
	Level  string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Format string `yaml:"format" env:"LOG_FORMAT" env-default:"text"`
}

func Load() *Config {
	_ = godotenv.Load(".env")
	var cfg Config
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "config/local.yaml"
	}
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("read env: %v", err)
	}
	return &cfg
}
