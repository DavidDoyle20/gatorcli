package config

import (
	"encoding/json"
	"os"
	"path"
)

var configFile = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}

func SetUser(cfg Config, username string) error {
	cfg.CurrentUserName = username
	err := write(cfg)
	if err != nil {
		return err
	}
	return nil
}

func SetDB(cfg Config, db_url string) error {
	cfg.DbURL = db_url
	err := write(cfg)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(home, configFile), nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func CreateConfigFile() error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	_, err = os.Stat(path)
	// file exists. do nothing
	if err == nil {
		return nil
	}
	_, err = os.Create(path)
	if err != nil {
		return err
	}
	return nil
}
