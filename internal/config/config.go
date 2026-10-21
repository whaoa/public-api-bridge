package config

import (
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/viper"
	"github.com/whaoa/public-api-bridge/pkg/errors"
)

const (
	configName      = "config"
	configNameLocal = configName + ".local"
)

type Config struct {
	LogLevel   string   `mapstructure:"log_level"`
	ServerPort int      `mapstructure:"server_port"`
	TMDBKeys   []string `mapstructure:"tmdb_keys"`
}

func Load() (*Config, error) {
	v := newViper()

	v.SetConfigName(configName)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.Wrap(err, "read config error")
		}
	}

	v.SetConfigName(configNameLocal)
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.Wrap(err, "read local config error")
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(err, "unmarshal config error")
	}

	config.TMDBKeys = lo.Uniq(lo.Filter(config.TMDBKeys, func(key string, _ int) bool {
		return strings.TrimSpace(key) != ""
	}))

	return &config, nil
}

// newViper builds an isolated loader so repeated Load calls do not accumulate
// state and each call resolves the working directory afresh.
func newViper() *viper.Viper {
	v := viper.New()

	v.SetConfigType("toml")
	v.AddConfigPath(".")
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	// Every key needs a default: AutomaticEnv only contributes to Unmarshal for
	// keys viper already knows about.
	v.SetDefault("log_level", "info")
	v.SetDefault("server_port", 3000)
	v.SetDefault("tmdb_keys", []string{})

	return v
}
