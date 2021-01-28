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
	cfg := filepath.Join(cfd, "ghf", "config.yaml")

	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		return err
	}

	viper.AddConfigPath(cfg)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if !viper.IsSet("token") {
		return errors.New("not found token")
	}

	viper.SetDefault("user", "unknown")
	viper.SetDefault("email", "unknown")

	return nil
}
