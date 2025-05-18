package config

import (
	"os"

	"github.com/joho/godotenv"
)

// config loads the application's configuration from environment variables

type Config struct {
	WebServer WebServerConfig
	DB        DBConfig
	TestDB    DBConfig
}

type WebServerConfig struct {
	Port string
	Host string
}

type DBConfig struct {
	DSN string
}

func (c *Config) LoadConfig() error {
	var err error
	if err = godotenv.Load(); err != nil {
		return ErrNoEnvFileFound
	}

	var envMap = map[string]*string{
		"DB_DSN":      &c.DB.DSN,
		"TEST_DB_DSN": &c.TestDB.DSN,
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
