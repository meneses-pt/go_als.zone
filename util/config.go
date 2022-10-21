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
}

func init() {
	loadAppConfig()
}

var AppConfig *Config

// loadAppConfig reads application configuration from environment variables.
func loadAppConfig() {
	log.Println("Loading Server Configurations...")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(AppConfig.DBHost)
}
