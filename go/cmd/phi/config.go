package main

import (
	"io"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/clstb/phi/go/pkg/config"
	"gopkg.in/yaml.v2"
)

func loadConfig(path string) tea.Cmd {
	return func() tea.Msg {
		config := config.Config{}

		f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(&config); err != nil {
			if err == io.EOF {
				return config
			}
			return err
		}

		return config
	}
}

func saveConfig(
	path string,
	config config.Config,
) tea.Cmd {
	return func() tea.Msg {
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}

		if err := yaml.NewEncoder(f).Encode(&config); err != nil {
			return err
		}

		return nil
	}
}
