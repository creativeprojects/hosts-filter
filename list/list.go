package list

import (
	"bufio"
	"io"
	"strings"

	"github.com/creativeprojects/clog"
	"github.com/creativeprojects/hosts-filter/constants"
)

func LoadLines(r io.Reader) ([]string, error) {
	lines := make([]string, 0, constants.BufferInitialLines)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, trim(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

func LoadEntries(lines []string, entries map[string]bool) {
	if entries == nil {
		entries = make(map[string]bool, constants.BufferInitialEntries)
	}
	for num, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		// remove any comment at the end of the line
		commentedOut := strings.Split(line, "#")
		parts := strings.Split(trim(commentedOut[0]), " ")
		if len(parts) != 2 {
			clog.Debugf("entry ignored one line %d: %q", num+1, line)
			continue
		}
		// check the IP is for a filtered domain
		if parts[0] != "0.0.0.0" && parts[0] != "127.0.0.1" {
			clog.Debugf("entry ignored one line %d: %q", num+1, line)
			continue
		}
		// check the entry is not ignored
		if isIgnored(parts[1]) {
			clog.Debugf("entry ignored one line %d: %q", num+1, line)
			continue
		}
		entries[parts[1]] = true
	}
}

func trim(value string) string {
	return strings.Trim(value, " \t\r\n")
}

func isIgnored(domain string) bool {
	for _, ignore := range constants.IgnoreDomains {
		if domain == ignore {
			return true
		}
	}
	return false
}
