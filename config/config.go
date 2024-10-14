package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Config structure to map the YAML fields
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Storage StorageConfig `yaml:"storage"`
	Debug   bool          `yaml:"debug"`
}

// ServerConfig for the server settings
type ServerConfig struct {
	Host           string   `yaml:"host"`
	Port           int      `yaml:"port"`
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

var Configuration *Config

func LoadConfig() {
	exist, err := DoesExistFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if !exist {
		err := writeConfig(Config{
			Server: ServerConfig{
				Host:           "127.0.0.1",
				Port:           8080,
				ServerURL:      "http://127.0.0.1:8080",
				TrustedProxies: []string{},
			},
			JWT: JWTConfig{
				Secret: "ASuperSecretSecretlyHiddenThatNobodyKnows",
			},
			Argon: ArgonConfig{
				Salt:        "ASuperSaltForArgon2idHashFunction",
				Parallelism: 2,
				Memory:      64,
				Iterations:  2,
				HashLenght:  32,
			},
			Database: DatabaseConfig{
				ServerURL:    "mongodb://mongouser:mongopass@localhost:27017",
				DatabaseName: "nebulogo",
			},
			Storage: StorageConfig{
				Directory: "storage",
			},
			Debug: false,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// Read YAML file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}

	// Unmarshal YAML data into Config struct
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshalling YAML: %v", err)
	}

	Configuration = &config
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

func DoesExistFile(filename string) (bool, error) {
	if _, err := os.Stat(filename); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}
