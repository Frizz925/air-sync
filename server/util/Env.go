package util

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

func GetEnvDefault(name string, def string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}
	return def
}

func EnvMongoUri() string {
	return GetEnvDefault("MONGODB_URI", "mongodb://localhost:27017")
}

func EnvMongoUrl() (*url.URL, error) {
	return url.Parse(EnvMongoUri())
}

func SetupCredentialsEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", path.Join(dir, "credentials.json"))
	return nil
}

func LoadDotEnv() error {
	mode := GetEnvDefault("APP_MODE", "development")
	filename := fmt.Sprintf(".env.%s", mode)
	var filenames []string
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		filenames = []string{".env"}
	} else if err == nil {
		filenames = []string{filename, ".env"}
	} else {
		return nil
	}
	log.Infof("Loading environment files: %s", strings.Join(filenames, ", "))
	return godotenv.Load(filenames...)
}
