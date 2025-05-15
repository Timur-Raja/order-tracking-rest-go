package app

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WebServer WebServerConfig
	DB        DBConfig
}

type WebServerConfig struct {
	Port string
	Host string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Username string
}

func (c *Config) LoadConfig() error {
	// Load environment variables from .env file
	var err error
	err = godotenv.Load()
	if err != nil {
		return ErrNoEnvFileFound
	}

	if c.DB.Host, err = getEnv("DB_HOST"); err != nil {
		return err
	}
	if c.DB.Port, err = getEnv("DB_PORT"); err != nil {
		return err
	}
	if c.DB.User, err = getEnv("DB_USER"); err != nil {
		return err
	}
	if c.DB.Password, err = getEnv("DB_PASSWORD"); err != nil {
		return err
	}
	if c.DB.Name, err = getEnv("DB_NAME"); err != nil {
		return err
	}
	if c.DB.Username, err = getEnv("DB_USERNAME"); err != nil {
		return err
	}
	if c.WebServer.Port, err = getEnv("APP_PORT"); err != nil {
		return err
	}
	if c.WebServer.Host, err = getEnv("APP_HOST"); err != nil {
		return err
	}
	return nil
}

func getEnv(key string) (string, error) {
	if key == "" {
		return "", ErrNoEnvKeyProvided
	}

	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	} else {
		return "", ErrNoEnvValueFound
	}
}
