package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	errRead = errors.New("config is empty, no config could be loaded")
)

// Emitterio struct contains this brokers configuration.
type Emitterio struct {
	Broker           string `json:"broker" env:"EMITTERIO_BROKER"`
	EventChannelKey  string `json:"eventChannelKey" env:"EMITTERIO_EVENT_CHAN_KEY"`
	EventChannelName string `json:"eventChannelName" env:"EMITTERIO_EVENT_CHAN_NAME"`
}

// Rest struct contains the endpoints' address and the path of the certification files.
type Rest struct {
	Endpoint string `json:"endpointTLS" env:"REST_ENDPOINT_TLS"`
	CertPath string `json:"certPath" env:"REST_CERT_PATH"`
}

// Prometheus struct represent the Prometheus' config.
type Prometheus struct {
	Endpoint string `json:"endpoint" env:"PROM_ENDPOINT"`
}

// Type represents the storage type
type StorageType string

// Supported types
const (
	StorageTypeMongo  = StorageType("mongo")
	StorageTypeDynamo = StorageType("dynamo")
)

type DB struct {
	Type       StorageType `json:"type" env:"DB_TYPE"`
	Connection string      `json:"connection" env:"DB_CONNECTION"`
	Region     string      `json:"region" env:"DB_REGION"`
}

// Config is the specific config to this service.
type Config struct {
	Rest       Rest       `json:"rest"`
	DB         DB         `json:"db"`
	Emitterio  Emitterio  `json:"emitterio"`
	Prometheus Prometheus `json:"prometheus"`
}

// Load the config from the file if any.
// Replace values read from the file with environment variable when found.
func Load(filename string) (*Config, error) {
	var cfg Config

	// Attempts to read the file.
	if file, err := os.Open(filename); err != nil {
		fmt.Println("Config: unable to read the config file. ", err.Error())
	} else {
		if err := json.NewDecoder(file).Decode(&cfg); err != nil {
			return nil, err
		}
	}

	// Read the env and replace variables found in the config file with the env version.
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		fmt.Println("Config: unable to get the config from the env. ", err.Error())
	}

	if cfg == (Config{}) {
		return nil, errRead
	}

	return &cfg, nil
}
