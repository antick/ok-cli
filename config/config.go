package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	// DefaultDestination is the default destination for copy and move operations
	DefaultDestination string `mapstructure:"default_destination"`
	// DefaultBuildOutput is the default output name for build operations
	DefaultBuildOutput string `mapstructure:"default_build_output"`
	// VerboseOutput enables more detailed output
	VerboseOutput bool `mapstructure:"verbose_output"`
	// PermanentDelete sets whether to permanently delete files by default
	PermanentDelete bool `mapstructure:"permanent_delete"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.ok")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			err = createDefaultConfig()
			if err != nil {
				return
			}
			err = viper.ReadInConfig()
			if err != nil {
				return
			}
		} else {
			// Config file was found but another error was produced
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}

func createDefaultConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %w", err)
	}

	configDir := filepath.Join(home, ".ok")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}
	defer f.Close()

	defaultConfig := `# Default destination for copy and move operations
default_destination: ""

# Default output name for build operations
default_build_output: "main"

# Enable verbose output
verbose_output: false

# Permanently delete files instead of moving to trash
permanent_delete: false
`

	_, err = f.WriteString(defaultConfig)
	if err != nil {
		return fmt.Errorf("error writing default config: %w", err)
	}

	return nil
}
