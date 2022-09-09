package config

import (
	"github.com/spf13/viper"
	"time"
)

const envFileName = ".env"

type Config struct{}

func NewConfig() *Config {
	viper.SetConfigFile(envFileName)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	return &Config{}
}

func (c *Config) GetGrpcPort() string {
	return viper.GetString("SELF_GRPC_PORT")
}

func (c *Config) GetDbHost() string {
	return viper.GetString("DB_HOST")
}

func (c *Config) GetDbUsername() string {
	return viper.GetString("DB_USER")
}

func (c *Config) GetDbName() string {
	return viper.GetString("DB_NAME")
}

func (c *Config) GetDbPort() string {
	return viper.GetString("DB_PORT")
}

func (c *Config) GetDbPassword() string {
	return viper.GetString("DB_PASSWORD")
}

func (c *Config) GetPwdSalt() string {
	return viper.GetString("PASSWORD_SALT")
}

func (c *Config) GetCodeSalt() string {
	return viper.GetString("CODE_SALT")
}

func (c *Config) GetJwtAccessSecretKey() string {
	return viper.GetString("JWT_ACCESS_SECRET_KEY")
}

func (c *Config) GetJwtRefreshSecretKey() string {
	return viper.GetString("JWT_REFRESH_SECRET_KEY")
}

func (c *Config) GetAccessDuration() time.Duration {
	duration, err := time.ParseDuration(viper.GetString("JWT_ACCESS_DURATION"))
	if err != nil {
		panic(err)
	}

	return duration
}

func (c *Config) GetRefreshDuration() time.Duration {
	duration, err := time.ParseDuration(viper.GetString("JWT_REFRESH_DURATION"))
	if err != nil {
		panic(err)
	}

	return duration
}

func (c *Config) GetNotifyGrpcPort() string {
	return viper.GetString("NOTIFY_SERVICE_GRPC_PORT")
}

func (c *Config) GetAppIssue() string {
	return viper.GetString("APP_ISSUE")
}

func (c *Config) GetEncryptKey() string {
	return viper.GetString("ENCRYPT_KEY")
}

func (c *Config) GetNotifyAddress() string {
	return viper.GetString("NOTIFY_ADDRESS")
}

func (c *Config) GetCodeLength() int {
	return viper.GetInt("CODE_SYMBOL_COUNT")
}
