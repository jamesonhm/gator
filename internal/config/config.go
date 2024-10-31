package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const config_filename = ".gatorconfig.json"

type Config struct {
	DB_url   string `json:"db_url"`
	CurrUser string `json:"current_user_name"`
}

func Read() Config {
	res := Config{}
	fp, err := getConfigFilepath()
	if err != nil {
		return res
	}
	byteData, _ := os.ReadFile(fp)
	err = json.Unmarshal(byteData, &res)
	if err != nil {
		fmt.Println("filepath:", fp)
		fmt.Println(res)
		return res
	}
	return res
}

func (c *Config) SetUser(user string) error {
	c.CurrUser = user
	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fp, _ := getConfigFilepath()
	err = os.WriteFile(fp, b, 0666)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilepath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("dir:", dir)
		fmt.Println("Error getting user dir")
		return "", err
	}
	filepath := dir + "/" + config_filename
	return filepath, nil
}
