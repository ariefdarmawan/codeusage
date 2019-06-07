package main

import (
	"github.com/spf13/viper"
)

// Config holds information about application configuration. It will reads from config file that will be passed to app using -config=path_to_config parameter
type Config struct {
	WorkingDir   string
	Libraries    []string
	Projects     []string
	Multiplier   int
	UsageKeyword []string
}

func readConfig(configPath string) (*Config, error) {
	cfg := new(Config)

	v := viper.New()
	v.SetConfigName("app.conf")
	v.AddConfigPath(configPath)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if err = v.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
