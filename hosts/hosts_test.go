package hosts

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	eol              string
	simpleContent    string
	ownContent       string
	generatedContent string
)

func init() {
	eol = "\n"
	simpleContent = "#\n#\n#\n127.0.0.1 localhost\n"
	ownContent = "-- 1\n#\n2\n3\n# --\n"
	generatedContent = "ip line1\nip line2\n"

	if runtime.GOOS == "windows" {
		eol = "\r\n"
		simpleContent = "#\r\n#\r\n#\r\n127.0.0.1 localhost\r\n"
		ownContent = "-- 1\r\n#\r\n2\r\n3\n# --\r\n"
		generatedContent = "ip line1\r\nip line2\r\n"
	}
}

func TestVirginHosts(t *testing.T) {
	contents := simpleContent
	result, _, found := extractOwnSection(contents)
	assert.False(t, found)
	assert.Equal(t, contents, result)
}

func TestOwnSection(t *testing.T) {
	own := ownContent
	before := simpleContent
	after := "# blah blah" + eol
	contents := before + startMarker + own + endMarker + after
	beforeResult, afterResult, found := extractOwnSection(contents)
	assert.True(t, found)
	assert.Equal(t, before, beforeResult)
	assert.Equal(t, after, afterResult)
}

func TestSectionOnItsOwn(t *testing.T) {
	own := ownContent
	contents := startMarker + own + endMarker
	beforeResult, afterResult, found := extractOwnSection(contents)
	assert.True(t, found)
	assert.Equal(t, "", beforeResult)
	assert.Equal(t, "", afterResult)
}

func TestUpdateEmptyHostsfile(t *testing.T) {
	buffer := &strings.Builder{}
	err := Update("", "ip", nil, buffer)
	require.NoError(t, err)
	assert.Equal(t, "", buffer.String())
}

func TestUpdateSimpleHostsfile(t *testing.T) {
	buffer := &strings.Builder{}
	err := Update("", "ip", []string{"line1", "line2"}, buffer)
	require.NoError(t, err)
	assert.Equal(t, eol+startMarker+generatedContent+endMarker, buffer.String())
}

func TestUpdateExistingHostsfile(t *testing.T) {
	buffer := &strings.Builder{}
	err := Update("something"+eol+startMarker+endMarker, "ip", []string{"line1", "line2"}, buffer)
	require.NoError(t, err)
	assert.Equal(t, "something"+eol+startMarker+generatedContent+endMarker, buffer.String())
}

func TestRemoveEntriesFromHostsfile(t *testing.T) {
	buffer := &strings.Builder{}
	err := Update(eol+startMarker+endMarker, "ip", nil, buffer)
	require.NoError(t, err)
	assert.Equal(t, eol, buffer.String())
}
