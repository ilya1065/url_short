package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	User        string        `yaml:"user"`
	Password    string        `yaml:"password" env:"HTTP_SERVER_PASSWORD" env-required:"true"`
}

func MustLoad() Config {
	configPath := *flag.String("config-path", "", "config-path")
	flag.Parse()
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("по указанному пути конфига нету путь: %s", configPath)

	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("не удалось считатать конфиг: %s", err)
	}

	return cfg
}
