package main

import (
	"log"
	"os"

	"github.com/creativeprojects/clog"
)

func setupConsoleLogger(flags commandLineFlags) {
	consoleHandler := clog.NewConsoleHandler("", 0)
	if flags.noAnsi {
		consoleHandler.Colouring(false)
	}
	logger := newFilteredLogger(flags, consoleHandler)
	clog.SetDefaultLogger(logger)
}

func setupFileLogger(flags commandLineFlags) (*os.File, error) {
	file, err := os.OpenFile(flags.logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	logger := newFilteredLogger(flags, clog.NewStandardLogHandler(file, "", log.LstdFlags))
	// default logger added with level filtering
	clog.SetDefaultLogger(logger)
	// and return the file handle (so we can close it at the end)
	return file, nil
}

func newFilteredLogger(flags commandLineFlags, handler clog.Handler) *clog.Logger {
	minLevel := clog.LevelInfo
	if flags.quiet {
		minLevel = clog.LevelWarning
	} else if flags.verbose {
		minLevel = clog.LevelDebug
	}
	// now create and return the logger
	return clog.NewLogger(clog.NewLevelFilter(minLevel, handler))
}

func changeLevelFilter(level clog.LogLevel) {
	handler := clog.GetDefaultLogger().GetHandler()
	filter, ok := handler.(*clog.LevelFilter)
	if ok {
		filter.SetLevel(level)
	}
}
