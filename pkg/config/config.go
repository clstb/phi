package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AccessToken string `yaml:"access_token"`
}

func Load(path string) (*Config, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Save(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	if err := yaml.NewEncoder(f).Encode(c); err != nil {
		return err
	}

	return nil
}
