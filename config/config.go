package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	App         AppConfig         `toml:"app"`
	YouTube     YouTubeConfig     `toml:"youtube"`
	Archive     ArchiveConfig     `toml:"archive"`
	CronPresets map[string]string `toml:"cron_presets"`
}

type AppConfig struct {
	MaxConcurrency int64 `toml:"max_concurrency"`
}

type YouTubeConfig struct {
	ApiKey string `toml:"api_key"`
}

type ArchiveConfig struct {
	ChannelsListURL string `toml:"channels_list_url"`
	ArchiveURL      string `toml:"archive_url"`
	QueueURL        string `toml:"queue_url"`
}

// LoadConfig tries to load the configuration from the file specified in the
// KUMO_CONFIG_FILE environment variable. If the variable is not set, it
// defaults to "config.toml".
func LoadConfig() (*Config, error) {
	configFileName := os.Getenv("KUMO_CONFIG_FILE")
	if configFileName == "" {
		configFileName = "config.toml"
	}

	var config Config
	if _, err := toml.DecodeFile(configFileName, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
