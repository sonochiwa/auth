package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type WebServerConfig struct {
	Port int
}

type DBConnectionConfig struct {
	Type     string
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Pool     struct {
		MaxIdleConns int `mapstructure:"max-idle-conns"`
		MaxOpenConns int `mapstructure:"max-open-conns"`
		IdleTimeout  int `mapstructure:"idle-timeout"`
	}
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
	DB    DBConnectionConfig      `mapstructure:"db"`
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

	viper.SetDefault("db.type", "postgres")
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.username", "root")
	viper.SetDefault("db.password", "root")
	viper.SetDefault("db.database", "auth")
	viper.SetDefault("db.pool.max-idle-conns", 10)
	viper.SetDefault("db.pool.max-open-conns", 10)
	viper.SetDefault("db.pool.max-idle-timeout", 300)

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
