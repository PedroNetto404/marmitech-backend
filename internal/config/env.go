package config

import (
	"fmt"
	"os"

	"github.com/codingconcepts/env"
	"github.com/joho/godotenv"
)

type environtment struct {
	ApiHost           string `env:"API_HOST"`
	ApiPort           string `env:"API_PORT" default:"8080"`
	ApiRouteTimeout   string `env:"API_ROUTE_TIMEOUT" default:"10000"`
	ApiBasePath       string `env:"API_BASE_PATH" default:"/api/v1"`
	DocEnabled        string `env:"DOC_ENABLED"`
	AuthEnabled       string `env:"AUTH_ENABLED"`
	KeycloakIssuerUrl string `env:"KEYCLOAK_ISSUER_URL"`
	KeycloakClientId  string `env:"KEYCLOAK_CLIENT_ID"`
	MySqlHost         string `env:"MYSQL_HOST"`
	MySqlPort         string `env:"MYSQL_PORT"`
	MySqlDatabase     string `env:"MYSQL_DATABASE"`
	MySqlUser         string `env:"MYSQL_USER"`
	MySqlPass         string `env:"MYSQL_PASS"`
	DiskStoragePath string `env:"DISK_STORAGE_PATH"`
	JwtSecretKey string `env:"JWT_SECRET_KEY"`
	JwtIssuer string `env:"JWT_ISSUER"`
	JwtAudience string `env:"JWT_AUDIENCE"`
	JwtExpirationMinutes int `env:"JWT_EXPIRATION_MINUTES"`
}

var Env environtment

func (e *environtment) IsDocEnabled() bool {
	return e.DocEnabled == "true"
}

func (e *environtment) IsAuthEnabled() bool {
	return e.AuthEnabled == "true"
}

func LoadEnvs() {
	goEnv := os.Getenv("ENV")
	if goEnv == "" {
		goEnv = "development"
	}

	var envPath string

	if os.Getenv("DEBUG") == "true" {
		envPath = fmt.Sprintf("../../.env.%s", goEnv)
	} else {
		envPath = fmt.Sprintf(".env.%s", goEnv)
	}

	if err := godotenv.Load(envPath); err != nil {
		panic(fmt.Errorf("failed to load env file: %v", err))
	}

	if err := env.Set(&Env); err != nil {
		panic(fmt.Errorf("failed to parse env: %v", err))
	}
}
