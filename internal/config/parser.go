package config

import (
	"setupx/internal/models"
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parse(data)
}

func LoadConfigFromURL(url string) (*models.Config, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch config: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parse(data)
}

func parse(data []byte) (*models.Config, error) {
	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
