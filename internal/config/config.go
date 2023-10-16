package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type WebServerConfig struct {
	Port int
}

type MongoDBConnectionConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type Config struct {
	Mongo MongoDBConnectionConfig `mapstructure:"mongodb"`
	Web   WebServerConfig         `mapstructure:"web"`
}

func InitConfiguration() *Config {
	var C *Config = new(Config)

	C.loadFile()
	C.loadDefault()

	viper.Unmarshal(C)
	return C
}

func (c *Config) loadDefault() {
	viper.SetDefault("mongodb.host", "localhost")
	viper.SetDefault("mongodb.port", 27017)
	viper.SetDefault("mongodb.database", "auth")
	viper.SetDefault("mongodb.username", "root")
	viper.SetDefault("mongodb.password", "root")

	viper.SetDefault("web.port", 8080)
}

func (c *Config) loadFile() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error reading config file, %s\n", err))
	}
}
