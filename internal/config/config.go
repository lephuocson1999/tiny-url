package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port         int           `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	Redis struct {
		Host string        `yaml:"host"`
		Port int           `yaml:"port"`
		TTL  time.Duration `yaml:"ttl"`
	} `yaml:"redis"`
	URLShortener struct {
		BaseURL    string        `yaml:"base_url"`
		CodeLength int           `yaml:"code_length"`
		DefaultTTL time.Duration `yaml:"default_ttl"`
	} `yaml:"urlshortener"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
