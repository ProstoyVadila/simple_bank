package utils

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

var TimeFormat = "02.01.2006 15:04:05"

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	GinMode             string        `mapstructure:"GIN_MODE"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig loads the config file from the given path
func LoadConfig(path string) (config Config, err error) {
	// TODO: rewrite this logic to getting filename from flag
	var filename string
	if isProduction() {
		filename = "prod"
	} else {
		filename = "dev"
	}

	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// isProduction checks if the prod.env file exists in the current directory
func isProduction() bool {
	if _, err := os.Stat("prod.env"); err == nil {
		return true
	}
	return false
}
