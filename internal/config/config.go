package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the struct that holds the configuration values
type ServerConfig struct {
	Host  string
	Port  int
	Debug bool
}

type JWTConfig struct {
	SecretKey  string
	Algorithm  string
	Expiration int
}

type LogConfig struct {
	Level  string
	Format string
	Output string
}

type Config struct {
	Server    ServerConfig
	JWT       JWTConfig
	LogConfig LogConfig
}

// LoadConfig loads the configuration values from the environment variables
func LoadConfig(appConfig *Config) {
	file, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&appConfig)
	if err != nil {
		panic(err)
	}

	appConfig.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")
}

func (c *Config) GetJWTSecretKey() (*ecdsa.PrivateKey, error) {
	der, err := base64.StdEncoding.DecodeString(c.JWT.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("500: Error decoding secret key")
	}

	privateKey, err := x509.ParseECPrivateKey(der)
	if err != nil {
		return nil, fmt.Errorf("500: Error parsing EC private key")
	}

	return privateKey, nil
}
