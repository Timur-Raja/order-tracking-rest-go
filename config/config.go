package config

import (
	"os"

	"github.com/joho/godotenv"
)

// config loads the application's configuration from environment variables

type Config struct {
	WebServer WebServerConfig
	DB        DBConfig
	ES        ESConfig
	Redis     RedisConfig
}

type WebServerConfig struct {
	URL string
}

type DBConfig struct {
	DSN     string
	TestDSN string
}

type ESConfig struct {
	URL     string
	TestURL string
}

type RedisConfig struct {
	URL     string
	TestURL string
}

func (c *Config) LoadConfig() error {
	var err error
	if err = godotenv.Load(); err != nil {
		return ErrNoEnvFileFound
	}

	var envMap = map[string]*string{
		"DB_DSN":         &c.DB.DSN,
		"TEST_DB_DSN":    &c.DB.TestDSN,
		"ES_URL":         &c.ES.URL,
		"TEST_ES_URL":    &c.ES.TestURL,
		"APP_URL":        &c.WebServer.URL,
		"REDIS_URL":      &c.Redis.URL,
		"TEST_REDIS_URL": &c.Redis.TestURL,
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
