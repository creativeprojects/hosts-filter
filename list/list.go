package list

import (
	"bufio"
	"io"
	"strings"

	"github.com/creativeprojects/hosts-filter/constants"
)

func LoadLines(r io.Reader) ([]string, error) {
	lines := make([]string, 0, constants.BUFFER_INITIAL_LINES)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, strings.Trim(scanner.Text(), " \t\r\n"))
	}

	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

func LoadEntries(lines []string, entries map[string]bool) {
	if entries == nil {
		entries = make(map[string]bool, constants.BUFFER_INITIAL_LINES)
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		entries[line] = true
	}
}
