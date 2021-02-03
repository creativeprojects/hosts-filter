package main

import (
	"os"

	"github.com/spf13/pflag"
)

type commandLineFlags struct {
	help    bool
	quiet   bool
	verbose bool
	config  string
	noAnsi  bool
	logFile string
}

// loadFlags loads command line flags
func loadFlags() (*pflag.FlagSet, commandLineFlags) {
	flagset := pflag.NewFlagSet("hosts-filter", pflag.ExitOnError)

	flags := commandLineFlags{}

	flagset.BoolVarP(&flags.help, "help", "h", false, "display this help")
	flagset.BoolVarP(&flags.quiet, "quiet", "q", false, "display only warnings and errors")
	flagset.BoolVarP(&flags.verbose, "verbose", "v", false, "display some debugging information")
	flagset.StringVarP(&flags.config, "config", "c", "config.yaml", "configuration file")
	flagset.StringVarP(&flags.logFile, "log", "l", "", "logs into a file instead of the console")
	flagset.BoolVar(&flags.noAnsi, "no-ansi", false, "disable ansi control characters (disable console colouring)")

	_ = flagset.Parse(os.Args[1:])

	return flagset, flags
}
