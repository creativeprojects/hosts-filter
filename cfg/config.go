package cfg

import (
	"io"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v2"
)

const (
	XDGAppName = "hosts-filter"
)

// Config from the file
type Config struct {
	HostsFile  string   `yaml:"hosts_file"`
	IP         string   `yaml:"ip"`
	BlockLists []List   `yaml:"block_lists"`
	Allow      []string `yaml:"allow"`
	AllowFrom  string   `yaml:"allow_from"`
}

// List configuration
type List struct {
	URL string `yaml:"url"`
}

// LoadFile loads the configuration from the file
func LoadFile(fileName string) (Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	return loadConfig(file)
}

func loadConfig(reader io.Reader) (Config, error) {
	config := Config{}
	decoder := yaml.NewDecoder(reader)
	err := decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func FindConfigurationFile(configFile string) (string, error) {
	// 1. Simple case: current folder (or rooted path)
	if fileExists(configFile) {
		return configFile, nil
	}

	// 2. Next we try xdg as the "standard" for user configuration locations
	xdgFilename, err := xdg.SearchConfigFile(filepath.Join(XDGAppName, configFile))
	if err == nil {
		if fileExists(xdgFilename) {
			return xdgFilename, nil
		}
	}
	// Not found
	return "", err
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
