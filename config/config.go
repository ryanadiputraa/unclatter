package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	*Server      `mapstructure:"server"`
	*Postgres    `mapstructure:"postgres"`
	*GoogleOauth `mapstructure:"google_oauth"`
}

type Server struct {
	Port        string `mapstructure:"port"`
	FrontendURL string `mapstructure:"fe_url"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type GoogleOauth struct {
	RedirectURL  string `mapstructure:"redirect_url"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	State        string `mapstructure:"state"`
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
