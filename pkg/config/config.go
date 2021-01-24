package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

type FileConfig struct {
	Regex      string `yaml:"regex"`
	Date       int    `yaml:"date"`
	DateFormat string `yaml:"date_format"`
	Amount     int    `yaml:"amount"`
	Currency   int    `yaml:"currency"`
	Entity     int    `yaml:"entity"`
	Reference  int    `yaml:"reference"`
}

type Config struct {
	AccessToken string       `yaml:"access_token"`
	Files       []FileConfig `yaml:"files"`
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

func (c *Config) ForFile(name string) (FileConfig, error) {
	var fileConfig FileConfig
	found := false
	for _, fc := range c.Files {
		re, err := regexp.Compile(fc.Regex)
		if err != nil {
			return FileConfig{}, err
		}
		if re.MatchString(name) {
			fileConfig = fc
			found = true
		}
	}
	if !found {
		return FileConfig{}, fmt.Errorf("no config for: %s", name)
	}

	return fileConfig, nil
}
