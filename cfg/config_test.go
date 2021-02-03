package cfg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyConfiguration(t *testing.T) {
	config, err := loadConfig(bytes.NewReader([]byte("---")))
	assert.NoError(t, err)
	assert.Len(t, config.Lists, 0)
}

func TestSimpleConfiguration(t *testing.T) {
	content := `---
lists:
  # comment
  - url: http://something
  -
    url: http://something
`
	config, err := loadConfig(bytes.NewReader([]byte(content)))
	assert.NoError(t, err)
	assert.Len(t, config.Lists, 2)
}
