package util

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvDefault(name string, def string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}
	return def
}

func EnvPort() string {
	return GetEnvDefault("PORT", "8080")
}

func EnvMongoUri() string {
	return GetEnvDefault("MONGODB_URI", "mongodb://localhost:27017")
}

func EnvMongoUrl() (*url.URL, error) {
	return url.Parse(EnvMongoUri())
}

func EnvMongoDbName() string {
	return GetEnvDefault("MONGODB_DATABASE", "airsync")
}

func EnvRedisAddr() string {
	return GetEnvDefault("REDIS_ADDR", "localhost:6379")
}

func EnvRedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}

func LoadDotEnv() error {
	mode := GetEnvDefault("APP_MODE", "development")
	filename := fmt.Sprintf(".env.%s", mode)
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return godotenv.Load()
	} else if err != nil {
		return err
	}
	return godotenv.Load(filename, ".env")
}
