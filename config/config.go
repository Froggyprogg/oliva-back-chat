package config

import "github.com/spf13/viper"

type Config struct {
	GRPC
}

type GRPC struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		GRPC: GRPC{
			Port: viper.GetString("grpc.port"),
		},
	}
}

func InitConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
