// Package config provides a Viper-based loader that merges a JSON config file with
// environment variable overrides. All WanderPlan services use this package.
package config

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Load initialises Viper from a JSON config file and then allows environment variables
// (with the given envPrefix) to override any key.
//
// The envPrefix is uppercased and prefixed to env var names, e.g. prefix "WANDERPLAN"
// maps the env var WANDERPLAN_DB_HOST to the Viper key "db.host".
func Load(configName, configPath, envPrefix string, out any) error {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("json")
	v.AddConfigPath(configPath)
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		log.Warn().Err(err).Str("config", configName).Msg("config file not found, relying on environment variables")
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.Unmarshal(out); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}
	return nil
}
