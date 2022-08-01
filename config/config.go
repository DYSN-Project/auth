package config

import (
	"github.com/joho/godotenv"
	"os"
	"time"
)

const EnvName = ".env"

type Config struct{}

func NewConfig() *Config {
	if err := godotenv.Load(EnvName); err != nil {
		panic(err)
	}
	return &Config{}
}

func (e *Config) GetJwtAccessSecretKey() string {
	return os.Getenv("JWT_ACCESS_SECRET_KEY")
}

func (e *Config) GetJwtRefreshSecretKey() string {
	return os.Getenv("JWT_REFRESH_SECRET_KEY")
}

func (e *Config) GetJwtRegSecretKey() string {
	return os.Getenv("JWT_REG_SECRET_KEY")
}

func (e *Config) GetDbHost() string {
	return os.Getenv("DB_HOST")
}

func (e *Config) GetDbUsername() string {
	return os.Getenv("DB_USERNAME")
}

func (e *Config) GetDbName() string {
	return os.Getenv("DB_NAME")
}

func (e *Config) GetDbPort() string {
	return os.Getenv("DB_PORT")
}

func (e *Config) GetDbPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func (e *Config) GetGrpcPort() string {
	return os.Getenv("GRPC_PORT")
}

func (e *Config) GetRestPort() string {
	return os.Getenv("REST_PORT")
}

func (e *Config) GetServerMode() string {
	return os.Getenv("SERVER_MODE")
}

func (e *Config) GetEnvironment() string {
	return os.Getenv("ENVIRONMENT")
}

func (e *Config) GetPwdSalt() string {
	return os.Getenv("PASSWORD_SALT")
}

func (e *Config) GetRegisterDuration() time.Duration {
	duration := os.Getenv("JWT_REGISTER_DURATION")

	result, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	return result
}

func (e *Config) GetAccessDuration() time.Duration {
	duration := os.Getenv("JWT_ACCESS_DURATION")

	result, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	return result
}

func (e *Config) GetRefreshDuration() time.Duration {
	duration := os.Getenv("JWT_REFRESH_DURATION")

	result, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	return result
}
