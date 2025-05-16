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
	Password string
	Name     string
	Username string
}

func (c *Config) LoadConfig() error {
	// Load environment variables from .env file
	var err error
	if err = godotenv.Load(); err != nil {
		return ErrNoEnvFileFound
	}

	var envMap = map[string]*string{
		"DB_HOST":     &c.DB.Host,
		"DB_PORT":     &c.DB.Port,
		"DB_PASSWORD": &c.DB.Password,
		"DB_NAME":     &c.DB.Name,
		"DB_USERNAME": &c.DB.Username,
		"APP_PORT":    &c.WebServer.Port,
		"APP_HOST":    &c.WebServer.Host,
	}

	for key, value := range envMap {
		if *value, err = getEnv(key); err != nil {
			return err
		}
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
