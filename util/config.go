package util

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// Config stores the application configuration from environment variables
type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	HTTPAddr   string `mapstructure:"HTTP_ADDR"`
	MediaRoot  string `mapstructure:"MEDIA_ROOT"`
	RedditRoot string `mapstructure:"REDDIT_ROOT"`
}

// LoadAppConfig reads application configuration from environment variables.
func LoadAppConfig() (*Config, error) {
	log.Println("Loading Server Configurations...")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var appConfig *Config
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		return nil, err
	}
	fmt.Println(appConfig.DBHost)
	return appConfig, nil
}
