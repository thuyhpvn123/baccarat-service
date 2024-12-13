package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	DatabaseURL        string
	StorageDatabaseURL string
	APIAddress         string

	PrivateKey_     string
	MetaNodeVersion string

	NodeAddress           string
	NodeConnectionAddress string

	StorageAddress           string
	StorageConnectionAddress string

	BaccaratAddress    string
	BaccaratABIPath    string

	DnsLink_ string

	PathLevelDB   string
	API_PORT string
}

var Config *AppConfig

func LoadConfig(configFilePath string) (*AppConfig, error) {
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func (c *AppConfig) DnsLink() string {
	return c.DnsLink_
}
