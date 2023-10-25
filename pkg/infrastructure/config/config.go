package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/tonytcb/party-invite/pkg/domain"
)

type CorrelationIDKey string

const (
	dublinLocationConfig = "dublin"

	CorrelationIDKeyName CorrelationIDKey = "correlation_id"
)

type Config struct {
	AppName        string `mapstructure:"APP_NAME"`
	HTTPPort       int    `mapstructure:"HTTP_PORT"`
	BaseLocation   string `mapstructure:"BASE_LOCATION"`
	LocationNearTo int32  `mapstructure:"LOCATION_NEAR_TO"`
}

func (c *Config) IsValid() error {
	if c.AppName == "" {
		return errors.Errorf("undefined APP_NAME env var")
	}
	if c.HTTPPort <= 0 {
		return errors.Errorf("undefined or invalid HTTP_PORT env var")
	}
	if c.BaseLocation != dublinLocationConfig {
		return errors.Errorf("invalid BASE_LOCATION env var")
	}
	if c.LocationNearTo <= 0 {
		return errors.Errorf("undefined or invalid LOCATION_NEAR_TO env var")
	}

	return nil
}

func Load(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrapf(err, "error to read config, path: %s", path)
	}

	var config = &Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, errors.Wrapf(err, "error to unmarshal config, path: %s", path)
	}

	return config, nil
}

func (c *Config) GetBaseLocation() *domain.Coordinate {
	switch c.BaseLocation {
	case dublinLocationConfig:
		return domain.DublinLocation

	default:
		return &domain.Coordinate{}
	}
}
