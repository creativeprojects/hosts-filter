package main

import (
	"fmt"
	"os"

	"github.com/creativeprojects/clog"
	"github.com/creativeprojects/hosts-filter/cfg"
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

	// keep this one last if possible (so it will be first at the end)
	defer showPanicData()

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
	fmt.Printf("%+v", c)
}
