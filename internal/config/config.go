package config

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

type PostgresConfig struct {
}

type Config struct {
	Port     string   `json:"port"`
	Backends []string `json:"backends"`
}

func RunConfig(log *logrus.Logger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Не удалось загрузить .env файл: %v", err)
	}

	pathConfig := os.Getenv("PATH_CONFIG")
	if pathConfig == "" {
		log.Print("PATH_CONFIG не указан")
	}

	file, err := os.Open(pathConfig)
	defer file.Close()
	if err != nil {
		log.Printf("Ошибка при открытии файла конфига: %s", pathConfig)
		return nil, err
	}

	var cfg Config
	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		log.Errorf("Ошибка парсинга конфига %s", pathConfig)
		return nil, err
	}

	return &cfg, nil
}
