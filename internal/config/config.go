package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	filepath, err := getConfigFilePath()
	if err != nil {
		return Config{}, nil
	}
	dat, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, nil
	}
	cfg := Config{}
	err = json.Unmarshal(dat, &cfg)
	if err != nil {
		return Config{}, nil
	}
	return cfg, nil
}

func (cfg *Config) SetUser(userName string) error {
	if userName == "" {
		return errors.New("no user name was given")
	}
	cfg.CurrentUserName = userName
	return write(cfg)
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + configFileName, nil
}

func write(cfg *Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}
