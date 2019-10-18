package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config structure
type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
		SSL  bool   `yaml:"ssl"`
		CA   string `yaml:"ca"`
		KEY  string `yaml:"key"`
	}
	Database struct {
		Port     string `yaml:"port"`
		Host     string `yaml:"host"`
		UserName string `yaml:"username"`
		Password string `yaml:"password"`
		SSL      bool   `yaml:"ssl"`
	}
}

// ConfigFile is the file name.
const ConfigFile = "smeter.yml"

var defaultConfig Config

// Initialize default config structure.
func init() {
	defaultConfig = Config{}
	defaultConfig.Server.Port = "9000"
	defaultConfig.Server.Host = "localhost"
	defaultConfig.Server.SSL = true
	defaultConfig.Database.Port = "5432"
	defaultConfig.Database.Host = "localhost"
	defaultConfig.Database.SSL = true
}

// SaveConfig save structure to file.
func SaveConfig(c Config, filename string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}

// LoadConfig loads file into structure.
func LoadConfig(filename string) (Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return defaultConfig, err
	}

	var c Config
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return defaultConfig, err
	}

	return c, nil
}
