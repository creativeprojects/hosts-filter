package hosts

import (
	"io"
	"runtime"
	"strings"

	"github.com/creativeprojects/hosts-filter/constants"
)

var (
	lineEnding  = "\n"
	startMarker string
	endMarker   string
)

func init() {
	if runtime.GOOS == "windows" {
		lineEnding = "\r\n"
	}
	startMarker = constants.StartMarker + lineEnding
	endMarker = constants.EndMarker + lineEnding
}

func Update(source, ip string, entries []string, w io.StringWriter) error {
	var err error

	before, after, sectionFound := extractOwnSection(source)

	_, err = w.WriteString(before)
	if err != nil {
		return err
	}

	if !sectionFound && len(entries) > 0 {
		// add a new line at the end of the file before adding our stuff
		_, err = w.WriteString(lineEnding)
		if err != nil {
			return err
		}
	}

	if len(entries) > 0 {
		_, err = w.WriteString(startMarker)
		if err != nil {
			return err
		}
		err = Generate(ip, entries, w)
		if err != nil {
			return err
		}

		_, err = w.WriteString(endMarker)
		if err != nil {
			return err
		}
	}

	if sectionFound {
		_, err = w.WriteString(after)
		if err != nil {
			return err
		}
	}
	return nil
}

func Generate(ip string, entries []string, w io.StringWriter) error {
	var err error

	for _, entry := range entries {
		if entry == "" {
			continue
		}
		_, err = w.WriteString(ip)
		if err != nil {
			return err
		}
		_, err = w.WriteString(" ")
		if err != nil {
			return err
		}
		_, err = w.WriteString(entry)
		if err != nil {
			return err
		}
		_, err = w.WriteString(lineEnding)
		if err != nil {
			return err
		}
	}
	return nil
}

// extractOwnSection returns before our section, and after it (if found):
//  - It is not returning both start and end markers.
//  - If not found, it returns the whole file content in the first string value
func extractOwnSection(contents string) (string, string, bool) {
	start := strings.Index(contents, startMarker)
	end := strings.Index(contents, endMarker)
	if start == -1 || end == -1 {
		return contents, "", false
	}
	return contents[:start], contents[end+len(endMarker):], true
}
