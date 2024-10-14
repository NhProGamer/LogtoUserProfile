package config

import (
	"crypto/rand"
	"encoding/hex"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Config structure to map the YAML fields
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Storage StorageConfig `yaml:"storage"`
	Logto   LogtoConfig   `yaml:"logto"`
	Debug   bool          `yaml:"debug"`
}

// ServerConfig for the server settings
type ServerConfig struct {
	Host           string   `yaml:"host"`
	Port           int      `yaml:"port"`
	Secret         string   `yaml:"secret"`
	ServerURL      string   `yaml:"server_url"`
	TrustedProxies []string `yaml:"trusted_proxies"`
}

type LogtoConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AppId     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

type StorageConfig struct {
	Directory string `yaml:"directory"`
}

func LoadConfig() (Config, error) {
	exist, err := doesExistFile("config.yaml")
	if err != nil {
		return Config{}, err
	}
	if !exist {
		log.Println("No config file found ! Writing a config file...")
		secret, err := generateRandomSecret(256)
		if err != nil {
			return Config{}, err
		}
		err = writeConfig(Config{
			Server: ServerConfig{
				Host:           "localhost",
				Port:           5000,
				Secret:         secret,
				ServerURL:      "http://localhost:5000",
				TrustedProxies: []string{},
			},
			Storage: StorageConfig{
				Directory: "user-data",
			},
			Logto: LogtoConfig{
				Endpoint:  "",
				AppId:     "",
				AppSecret: "",
			},
			Debug: false,
		})
		if err != nil {
			return Config{}, err
		}
	}

	// Read YAML file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}

	// Unmarshal YAML data into Config struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func writeConfig(config Config) error {

	// Marshal struct into YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	// Write YAML data to file
	err = os.WriteFile("config.yaml", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func doesExistFile(file string) (bool, error) {
	if _, err := os.Stat(file); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func generateRandomSecret(length int) (string, error) {
	// Create a byte slice to hold the random data
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Convert the bytes to a hex string
	secret := hex.EncodeToString(bytes)
	return secret, nil
}
