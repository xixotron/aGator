package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cfgPath := filepath.Join(home, configFileName)

	return cfgPath, nil
}

func write(cfg Config) error {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(cfgPath)
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
