package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = "PATH_CONFIG"
)

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func check(path string) error {
	if path == "" {
		return fmt.Errorf("%s is not set", configPath)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", path)
	}

	return nil
}

func MustLoadWithDefault(def string) *Config {
	cfg, err := LoadWithDefault(def)
	if err != nil {
		panic(err)
	}
	return cfg
}

func LoadWithDefault(def string) (*Config, error) {
	const op = "config.Load"

	path := os.Getenv(configPath)

	if path == "" {
		path = def
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist", path)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return cfg, nil
}

func Load() (*Config, error) {
	const op = "config.Load"

	path := os.Getenv(configPath)

	if err := check(path); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return cfg, nil
}

func MustLoadWithEnv() *Config {
	cfg, err := LoadWithEnv()
	if err != nil {
		panic(err)
	}
	return cfg
}

func LoadWithEnv() (*Config, error) {
	const op = "config.LoadWithEnv"

	path := os.Getenv(configPath)
	if err := check(path); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return cfg, nil
}

func MustLoadLogger() *LoggerConfig {
	cfg, err := LoadLogger()
	if err != nil {
		panic(err)
	}
	return cfg
}

func LoadLogger() (*LoggerConfig, error) {
	const op = "config.LoadLogger"

	path := os.Getenv(configPath)
	if err := check(path); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	cfg := struct {
		Logger LoggerConfig `yaml:"logger"`
	}{}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return &cfg.Logger, nil
}

func LoadServer() (*ServerConfig, error) {
	const op = "config.LoadServer"

	path := os.Getenv(configPath)
	if err := check(path); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	cfg := struct {
		Server ServerConfig `yaml:"server"`
	}{}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return &cfg.Server, nil
}

func LoadServerWithEnv() (*ServerConfig, error) {
	const op = "config.LoadServer"

	path := os.Getenv(configPath)
	if err := check(path); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	cfg := struct {
		Server ServerConfig `yaml:"server"`
	}{}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	if err := cleanenv.ReadEnv(&cfg.Server); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return &cfg.Server, nil
}
