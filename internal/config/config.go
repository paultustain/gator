package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFile = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigPath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(path, configFile)

	return fullPath, nil
}

func Read() (Config, error) {
	fullPath, err := getConfigPath()

	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return Config{}, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)

	if err != nil {
		return Config{}, err
	}

	return cfg, nil

}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName

	return write(*cfg)
}

func write(cfg Config) error {
	fullPath, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
