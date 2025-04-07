package config

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

type env struct {
	ApiHost         string `env:"API_HOST"`
	ApiPort         string `env:"API_PORT" default:"8080"`
	ApiRouteTimeout string `env:"API_ROUTE_TIMEOUT" default:"10000"`
	DocEnabled      string `env:"DOC_ENABLED"`
	AuthEnabled     string `env:"AUTH_ENABLED"`

	MongoUrl string `env:"MONGO_URL"`
	MongoDb  string `env:"MONGO_DB"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     string `env:"REDIS_PORT"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       string `env:"REDIS_DB"`

	KeycloakIssuerUrl string `env:"KEYCLOAK_ISSUER_URL"`
	KeycloakClientId  string `env:"KEYCLOAK_CLIENT_ID"`
}

func (e *env) IsDocEnabled() bool {
	return e.DocEnabled == "true"
}

func (e *env) IsAuthEnabled() bool {
	return e.AuthEnabled == "true"
}

func LoadEnvs() {
	loadEnvFile()

	if err := loadEnvVars(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}
}

var Env env

func loadEnvFile() {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "development"
	}

	envPath := fmt.Sprintf(".env.%s", environment)

	file, err := os.Open(envPath)
	if err != nil {
		log.Fatal()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		keyValue := strings.SplitN(line, "=", 2)
		if len(keyValue) != 2 {
			log.Fatalf("invalid env line: %s", line)
		}
		key, value := keyValue[0], keyValue[1]

		if key == "" {
			log.Fatal("invalid env key")
		}

		os.Setenv(key, value)
	}
}

func loadEnvVars() error {
	v := reflect.ValueOf(&Env)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		value := os.Getenv(envTag)
		if value == "" {
			defaultValue := field.Tag.Get("default")
			if defaultValue != "" {
				value = defaultValue
			}
		}

		if value != "" {
			fieldValue := v.Field(i)
			if fieldValue.CanSet() {
				fieldValue.Set(reflect.ValueOf(value))
			}
		}
	}

	return nil
}
