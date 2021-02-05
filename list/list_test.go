package list

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortedKeys(t *testing.T) {
	testData := []map[string]bool{
		{"line1": true, "line2": true, "line3": true},
		{"line2": true, "line1": true, "line3": true},
		{"line3": true, "line1": true, "line2": true},
		{"line3": true, "line2": true, "line1": true},
	}
	for _, testMap := range testData {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, []string{"line1", "line2", "line3"}, SortedKeys(testMap))
		})
	}
}

func TestLoadLines(t *testing.T) {
	content := "\n1\n2\r\n3\n4\r\n5"
	buffer := bytes.NewBufferString(content)
	lines, err := LoadLines(buffer)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"", "1", "2", "3", "4", "5"}, lines)
}

func TestLoadEntries(t *testing.T) {
	content := `
# comment here
127.0.0.1 localhost
192.168.0.1 localhost
192.168.0.1 machine
127.0.0.1 filter1
0.0.0.0 filter2
nospacehere
`
	buffer := bytes.NewBufferString(content)
	lines, err := LoadLines(buffer)
	require.NoError(t, err)
	entries := make(map[string]bool, 10)
	LoadEntries(lines, entries)
	assert.ElementsMatch(t, []string{"filter1", "filter2"}, SortedKeys(entries))
}
