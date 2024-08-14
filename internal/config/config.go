package config

import (
	"github.com/VadimGossip/mm_agent/internal/model"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

func parseConfigFile(configDir string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func setFromEnv(cfg *model.Config) error {
	if err := envconfig.Process("db", &cfg.OracleDb); err != nil {
		return err
	}

	return envconfig.Process("mm", &cfg.Mattermost)
}

func unmarshal(cfg *model.Config) error {
	if err := viper.UnmarshalKey("oracle", &cfg.OracleDb); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("mm", &cfg.Mattermost); err != nil {
		return err
	}

	return nil
}

func Init(configDir string) (*model.Config, error) {
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}
	cfg := &model.Config{}
	if err := unmarshal(cfg); err != nil {
		return nil, err
	}
	if err := setFromEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
