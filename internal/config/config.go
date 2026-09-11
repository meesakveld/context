package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const FileName = ".context.yaml"

type Config struct {
	Output       string   `yaml:"output"`
	Format       string   `yaml:"format"`
	MaxFileSize  string   `yaml:"max_file_size"`
	IncludeEnv   bool     `yaml:"include_env"`
	TreeOnly     bool     `yaml:"tree_only"`
	FilesOnly    bool     `yaml:"files_only"`
	IncludeNotes bool     `yaml:"include_notes"`
	Exclude      []string `yaml:"exclude"`
	Include      []string `yaml:"include"`
}

func Default() Config {
	return Config{
		Output:       "context.txt",
		Format:       "txt",
		MaxFileSize:  "1MB",
		IncludeNotes: true,
	}
}

func Load() (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(FileName)

	if os.IsNotExist(err) {
		return cfg, nil
	}

	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
