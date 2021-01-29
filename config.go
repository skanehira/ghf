package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func initConfig() error {
	cfd, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	cfg := filepath.Join(cfd, "ghf")

	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		return err
	}

	viper.AddConfigPath(cfg)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if !viper.IsSet("token") {
		return errors.New("not found token")
	}

	if !viper.IsSet("user") {
		return errors.New("not found user")
	}

	viper.SetDefault("email", "unknown")

	return nil
}
