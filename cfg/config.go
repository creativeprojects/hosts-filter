package cfg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/adrg/xdg"
	"github.com/creativeprojects/hosts-filter/constants"
	"gopkg.in/yaml.v2"
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
	// 1. search from the current folder (or rooted path)
	if fileExists(configFile) {
		return configFile, nil
	}

	if filepath.IsAbs(configFile) {
		// no need to search further
		return "", fmt.Errorf("config file %q does not exist", configFile)
	}

	// 2. try xdg as the "standard" for user configuration locations
	xdgFilename, err := xdg.SearchConfigFile(filepath.Join(constants.Name, configFile))
	if err == nil {
		if fileExists(xdgFilename) {
			return xdgFilename, nil
		}
	}

	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		// Not found
		return "", err
	}

	// 3. search some standard unix paths
	for _, configPath := range constants.ConfigLocationsUnix {
		filename := filepath.Join(configPath, configFile)
		if fileExists(filename) {
			return filename, nil
		}
	}

	// Not found
	return "", fmt.Errorf(
		"could not locate %q in any of the following paths: %s",
		configFile,
		strings.Join(append(xdg.ConfigDirs, constants.ConfigLocationsUnix...), ", "),
	)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
