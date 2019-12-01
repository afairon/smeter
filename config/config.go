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
		DBName   string `yaml:"db"`
		SSL      bool   `yaml:"ssl"`
	}
}

// ConfigFile is the file name.
const ConfigFile = "smeter.yml"

var defaultConfig Config

// Initialize default config structure.
func init() {
	defaultConfig = Config{}
	defaultConfig.Server.Port = "50051"
	defaultConfig.Server.Host = "localhost"
	defaultConfig.Server.SSL = true
	defaultConfig.Database.Port = "5432"
	defaultConfig.Database.Host = "localhost"
	defaultConfig.Database.DBName = "smeter"
	defaultConfig.Database.SSL = true
}

// NewConfig returns the default config.
func NewConfig() *Config {
	return &defaultConfig
}

// Load reads config file and unmarshal it to the structure.
func (c *Config) Load(filename string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, c)
}

// Write takes the config structure, marshal and write it to file.
func (c *Config) Write(filename string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}
