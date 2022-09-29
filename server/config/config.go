package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database Database `mapstructure:"database"`
	JWTKey   string   `mapstructure:"jwt_key"`
	Server   Server   `mapstructure:"server"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Name     string `mapstructure:"db_name"`
	Password string `mapstructure:"password"`
	SslMode  string `mapstructure:"ssl_mode"`
}

type Server struct {
	Port uint `mapstructure:"port"`
}

func ReadConfig() (*Config, error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	//TODO throw errors if a required value is not set
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil

}
