package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

const configPath = "config/local.yaml"

type Config struct {
	Env             string `yaml:"env" env-default:"local"`
	FileStoragePath string `yaml:"file-storage-path" env-required:"true"`
	HTTPServer      `yaml:"http-server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" evn-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"60s"`
}

func MustLoad() Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println()
		log.Fatalf("Config file does not exists: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Unable to read config: %s", err)
	}
	return cfg
}
