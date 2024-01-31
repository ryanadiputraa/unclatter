package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	*Server
	*Postgres
}

type Server struct {
	Port string `mapstructure:"port"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func LoadConfig(configType, filePath string) (*Config, error) {
	viper.SetConfigType(configType)
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return config, nil
}
