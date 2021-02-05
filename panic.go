package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"text/tabwriter"

	"github.com/creativeprojects/hosts-filter/constants"
)

func showPanicData() {
	if r := recover(); r != nil {
		message := `
===============================================================
uh-oh! %s crashed miserably :-(
Can you please open an issue on github including these details:
===============================================================
`
		fmt.Fprintf(os.Stderr, message, constants.Name)
		w := tabwriter.NewWriter(os.Stderr, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "\t%s:\t%s\n", "os", runtime.GOOS)
		_, _ = fmt.Fprintf(w, "\t%s:\t%s\n", "arch", runtime.GOARCH)
		_, _ = fmt.Fprintf(w, "\t%s:\t%s\n", "version", version)
		_, _ = fmt.Fprintf(w, "\t%s:\t%s\n", "commit", commit)
		_, _ = fmt.Fprintf(w, "\t%s:\t%s\n", "compiled", date)
		_, _ = fmt.Fprintf(w, "\t%s:\t%s\n", "by", builtBy)
		_, _ = fmt.Fprintf(w, "\t%s:\t%s\n", "error", r)
		_, _ = fmt.Fprintf(w, "\t%s:\t%s\n", "stack", string(debug.Stack()))
		w.Flush()
		fmt.Fprint(os.Stderr, "===============================================================\n")
	}
}
