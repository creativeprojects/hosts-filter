package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedKeys(t *testing.T) {
	testData := []map[string]bool{
		{"line1": true, "line2": true, "line3": true},
		{"line2": true, "line1": true, "line3": true},
		{"line3": true, "line1": true, "line2": true},
		{"line3": true, "line2": true, "line1": true},
	}
	for _, testMap := range testData {
		assert.Equal(t, []string{"line1", "line2", "line3"}, sortedKeys(testMap))
	}
}
