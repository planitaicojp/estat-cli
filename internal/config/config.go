package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultFormat = "table"
	DefaultLang   = "J"
	configFile    = "config.yaml"
)

type Config struct {
	AppID  string `yaml:"app_id"`
	Format string `yaml:"format"`
	Lang   string `yaml:"lang"`
}

func ConfigDir() string {
	if d := os.Getenv(EnvConfigDir); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "estat")
}

func Load() (*Config, error) {
	path := filepath.Join(ConfigDir(), configFile)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルのパースに失敗: %w", err)
	}
	if cfg.Format == "" {
		cfg.Format = DefaultFormat
	}
	if cfg.Lang == "" {
		cfg.Lang = DefaultLang
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("設定のマーシャルに失敗: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, configFile), data, 0600)
}

func defaultConfig() *Config {
	return &Config{
		Format: DefaultFormat,
		Lang:   DefaultLang,
	}
}
