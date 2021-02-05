package main

import (
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/creativeprojects/clog"
	"github.com/creativeprojects/hosts-filter/cfg"
	"github.com/creativeprojects/hosts-filter/constants"
	"github.com/creativeprojects/hosts-filter/hosts"
	"github.com/creativeprojects/hosts-filter/list"
)

// These fields are populated by the goreleaser build
var (
	version = "0.1.0-dev"
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

	banner()

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
			c.HostsFile = constants.DefaultWindowsHostFile
		} else {
			c.HostsFile = constants.DefaultUnixHostFile
		}
	}
	c.HostsFile = expandEnv(c.HostsFile)

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
			if allow == "" {
				continue
			}
			delete(entries, allow)
		}

		if c.AllowFrom != "" {
			allowList, err := loadAllowList(c.AllowFrom)
			if err != nil {
				clog.Errorf("error while loading the allow list: %v", err)
			} else {
				for _, allow := range allowList {
					if allow == "" {
						continue
					}
					delete(entries, allow)
				}
			}
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

	if flags.stdout {
		// send the file to the default output
		err = hosts.Update(source, c.IP, list.SortedKeys(entries), os.Stdout)
		if err != nil {
			clog.Errorf("cannot write to standard output: %v", err)
			exitCode = 1
			return
		}
		return
	}

	// save the hosts file back (via a temporary file)
	tempfile := c.HostsFile + "-" + randSeq(6)
	err = saveHostsfile(source, c.IP, list.SortedKeys(entries), tempfile)
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
	clog.Debugf("file %q successfully saved", c.HostsFile)
}

func banner() {
	clog.Debugf("%s %s compiled with %s", constants.Name, version, runtime.Version())
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

func expandEnv(value string) string {
	if runtime.GOOS == "windows" {
		// variable expansion is always done the "unix" way, but never the "windows" way
		pattern := regexp.MustCompile(`%([^%]+)%`)
		value = pattern.ReplaceAllString(value, `${$1}`)
	}
	return os.ExpandEnv(value)
}

func loadAllowList(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return list.LoadLines(file)
}
