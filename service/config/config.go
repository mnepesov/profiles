package config

import (
	"github.com/spf13/viper"

	"github.com/mnepesov/profiles/service/server"
)

type Root struct {
	Server server.Config `mapstructure:"server"`
	Admin  Admin         `mapstructure:"admin"`
}

type Admin struct {
	Username string `mapstructure:"username"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

func NewConfig(configFile string) (*Root, error) {
	config, err := loadConfig(configFile)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func loadConfig(configFile string) (*Root, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var appConfig Root
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, err
}
