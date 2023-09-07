package config

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type Config struct {
	Database       DatabaseConfig `json:"database"`
	TokenPath      string         `json:"token_path"`
	Router         RouterConfig   `json:"router"`
	CredentialPath string         `json:"credential_path"`
	GoogleCOnfig   *oauth2.Config `json:"google_config,omitemppty"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type RouterConfig struct {
	Port         string `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

func ParseConfig(fileName string) (*Config, error) {
	var c Config
	cnfFile, err := os.ReadFile(path.Clean(fileName))
	if err != nil {
		log.Printf("unable to read config file: %v", err)
		return nil, err
	}

	if err = json.Unmarshal(cnfFile, &c); err != nil {
		log.Printf("unable to decode config file: %v", err)
		return nil, err
	}

	log.Println(c)
	credential, err := os.ReadFile(path.Clean(c.CredentialPath))
	if err != nil {
		log.Printf("unable to read credential file: %v", err)
		return nil, err
	}

	c.GoogleCOnfig, err = google.ConfigFromJSON(credential, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Printf("unable to get google config from json: %v", err)
		return nil, err
	}

	return &c, nil
}
