// File: internal/config/config.go

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	SuperOrganization struct {
		Name  string `yaml:"name"`
		Admin struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"admin"`
	} `yaml:"super_organization"`
	Crypto struct {
		Path string `yaml:"path"`
	} `yaml:"crypto"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}
