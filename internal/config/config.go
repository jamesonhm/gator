package config

import (
	"encoding/json"
	"os"
)

const config_filename = ".gatorconfig.json"

type Config struct {
	DBurl    string `json:"db_url"`
	CurrUser string `json:"current_user_name"`
}

func Read() (Config, error) {
	res := Config{}
	fp, err := getConfigFilepath()
	if err != nil {
		return res, err
	}
	file, err := os.Open(fp)
	if err != nil {
		return res, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrUser = user
	fp, err := getConfigFilepath()
	if err != nil {
		return err
	}

	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilepath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filepath := dir + "/" + config_filename
	return filepath, nil
}
