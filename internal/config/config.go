package config

import (
	"log"

	"time"

	"github.com/ilyakaznacheev/cleanenv"
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
	Address     string        `yaml:"address" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres123"`
	DBName   string `yaml:"dbname" env-default:"L0_db"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

type Kafka struct {
	Brokers  []string `yaml:"brokers" env:""`
	Topic    string   `yaml:"topic" env-default:"orders"`
	GroupID  string   `yaml:"group_id" env-default:"orderservice"`
	DLQTopic string   `yaml:"dlq_topic" env-default:"orders-dlq"`
	ClientID string   `yaml:"client_id" env-default:"orderservice"`
}

type Cache struct {
	MaxEntries int `yaml:"max_entries" env-default:"1000"`
}

type Log struct {
	Level  string `yaml:"level" env-default:"info"`
	Format string `yaml:"format" env-default:"text"`
}

func Load() *Config {
	configPath := "C:/Users/knyaz/Desktop/WB/ex-L0/config/local.yaml"
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	return &cfg
}
