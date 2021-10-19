package config

import (
	"io"
	"os"
	"path/filepath"

	"filippo.io/age"
	"github.com/clstb/phi/go/pkg/services/tinkgw/pb"
	"gopkg.in/yaml.v2"
)

type Config struct {
	TinkToken   *pb.Token
	AccessToken string
	Identity    string
	Recipient   string
}

func NewConfig() *Config {
	return &Config{
		TinkToken: &pb.Token{},
	}
}

func Load(path string) (*Config, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}

	c := NewConfig()
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		if err == io.EOF {
			return c, nil
		}
		return nil, err
	}

	if c.Identity == "" {
		identity, err := age.GenerateX25519Identity()
		if err != nil {
			return nil, err
		}

		recipient := identity.Recipient()

		c.Identity = identity.String()
		c.Recipient = recipient.String()
	}

	return c, nil
}

func (c *Config) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	if err := yaml.NewEncoder(f).Encode(c); err != nil {
		return err
	}

	return nil
}
