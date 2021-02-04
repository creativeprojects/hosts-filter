package main

import (
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/creativeprojects/clog"
	"github.com/creativeprojects/hosts-filter/cfg"
	"github.com/creativeprojects/hosts-filter/constants"
	"github.com/creativeprojects/hosts-filter/hosts"
	"github.com/creativeprojects/hosts-filter/list"
)

// These fields are populated by the goreleaser build
var (
	version = "0.11.1-dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	var exitCode = 0
	var err error

	// trick to run all defer functions before returning with an exit code
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	flagset, flags := loadFlags()
	// help
	if flags.help {
		flagset.Usage()
		return
	}

	if flags.logFile != "" {
		file, err := setupFileLogger(flags)
		if err != nil {
			// back to a console logger
			setupConsoleLogger(flags)
			clog.Errorf("cannot open logfile: %s", err)
		} else {
			// close the log file at the end
			defer file.Close()
		}

	} else {
		// Use the console logger
		setupConsoleLogger(flags)
	}

	// keep this defer last if possible (so it will be first at the end)
	defer showPanicData()

	// seed random number generation
	rand.Seed(time.Now().UnixNano())

	configFile, err := cfg.FindConfigurationFile(flags.config)
	if err != nil {
		clog.Error(err)
		exitCode = 1
		return
	}
	if configFile != flags.config {
		clog.Infof("using configuration file: %s", configFile)
	}

	c, err := cfg.LoadFile(configFile)
	if err != nil {
		clog.Errorf("cannot load configuration file: %v", err)
		exitCode = 1
		return
	}

	// validate missing configuration with defaults
	if c.IP == "" {
		c.IP = constants.DefaultFilteredIP
	}
	if c.HostsFile == "" {
		if runtime.GOOS == "windows" {
			c.HostsFile = os.ExpandEnv(constants.DefaultWindowsHostFile)
		} else {
			c.HostsFile = os.ExpandEnv(constants.DefaultUnixHostFile)
		}
	}

	var entries map[string]bool

	// load the entries if not in remove mode - otherwise empty entries will remove the section from the hosts file
	if !flags.remove {
		entries = make(map[string]bool, constants.BufferInitialEntries)
		for _, def := range c.BlockLists {
			if def.URL == "" {
				continue
			}
			err := loadListfile(def.URL, entries)
			if err != nil {
				clog.Error(err)
				continue
			}
		}

		if len(entries) == 0 {
			clog.Warning("the blocklist is empty")
			exitCode = 1
			return
		}

		// remove entries from the allow list
		for _, allow := range c.Allow {
			delete(entries, allow)
		}
	}

	var source string
	if fileExists(c.HostsFile) {
		source, err = loadHostsfile(c.HostsFile)
		if err != nil {
			clog.Errorf("cannot read hosts file: %v", err)
			exitCode = 1
			return
		}
	} else {
		clog.Warningf("hosts file %q does not exist and will be created.", c.HostsFile)
	}

	tempfile := c.HostsFile + "-" + randSeq(6)
	err = saveHostsfile(source, c.IP, sortedKeys(entries), tempfile)
	if err != nil {
		clog.Errorf("cannot write to temporary file: %v", err)
		exitCode = 1
		return
	}
	// now rename the temp file into the hosts file
	err = os.Rename(tempfile, c.HostsFile)
	if err != nil {
		clog.Errorf("cannot move the temporary file into place: %v", err)
		exitCode = 1
		return
	}
}

func loadListfile(filename string, entries map[string]bool) error {
	var reader io.Reader
	URL, err := url.Parse(filename)
	if err != nil || URL.Scheme == "" {
		// entry should be a file on disk
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		reader = file
	} else {
		// entry is http resource
		resp, err := http.Get(filename)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		reader = resp.Body
	}

	lines, err := list.LoadLines(reader)
	if err != nil {
		return err
	}

	list.LoadEntries(lines, entries)
	clog.Infof("loaded %q: %d entries in total", filename, len(entries))
	return nil
}

func sortedKeys(input map[string]bool) []string {
	output := make([]string, len(input))
	index := 0
	for key, _ := range input {
		output[index] = key
		index++
	}
	sort.Slice(output, func(i, j int) bool {
		return output[i] < output[j]
	})
	return output
}

func loadHostsfile(hostsfile string) (string, error) {
	file, err := os.Open(hostsfile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func saveHostsfile(source, ip string, entries []string, dest string) error {
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	return hosts.Update(source, ip, entries, file)
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
