package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string
}

const configFileName = ".tryConfig.json"

func ConfigRead() (*Config, error) {
	user, err := user.Current()
	if err != nil {
		fmt.Printf("%s", err)
	}
	filePath := filepath.Join(user.HomeDir, configFileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("%s", err)
	}
	var jsonData Config
	err = json.Unmarshal(content, &jsonData)
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}
	jsonData.CurrentUserName = "Vortex"

	//Just checking
	// fmt.Printf("%s\n", jsonData.DbUrl)
	// fmt.Printf("%s\n", jsonData.CurrentUserName)

	return &jsonData, nil
}
