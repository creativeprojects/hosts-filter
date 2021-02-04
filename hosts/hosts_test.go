package hosts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVirginHosts(t *testing.T) {
	contents := "#\n#\n#\n127.0.0.1 localhost\n"
	result, _, found := extractOwnSection(contents)
	assert.False(t, found)
	assert.Equal(t, contents, result)
}

func TestOwnSection(t *testing.T) {
	own := "-- 1\n#\n2\n3\n# --\n"
	before := "#\n#\n#\n127.0.0.1 localhost\n"
	after := "# blah blah\n"
	contents := before + startMarker + own + endMarker + after
	beforeResult, afterResult, found := extractOwnSection(contents)
	assert.True(t, found)
	assert.Equal(t, before, beforeResult)
	assert.Equal(t, after, afterResult)
}

func TestSectionOnItsOwn(t *testing.T) {
	own := "-- 1\n#\n2\n3\n# --\n"
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
	assert.Equal(t, "\n"+startMarker+endMarker, buffer.String())
}

func TestUpdateSimpleHostsfile(t *testing.T) {
	buffer := &strings.Builder{}
	err := Update("", "ip", []string{"line1", "line2"}, buffer)
	require.NoError(t, err)
	assert.Equal(t, "\n"+startMarker+"ip line1\nip line2\n"+endMarker, buffer.String())
}

func TestUpdateExistingHostsfile(t *testing.T) {
	buffer := &strings.Builder{}
	err := Update("something\n"+startMarker+endMarker, "ip", []string{"line1", "line2"}, buffer)
	require.NoError(t, err)
	assert.Equal(t, "something\n"+startMarker+"ip line1\nip line2\n"+endMarker, buffer.String())
}
