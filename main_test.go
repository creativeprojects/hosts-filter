package main

import (
	"os"
	"runtime"
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
		t.Run("", func(t *testing.T) {
			assert.Equal(t, []string{"line1", "line2", "line3"}, sortedKeys(testMap))
		})
	}
}

func TestExpandEnv(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Test designed for Unix systems")
	}
	testData := []struct {
		value    string
		expected string
	}{
		{"$HOME", os.Getenv("HOME")},
		{"${HOME}", os.Getenv("HOME")},
	}

	for _, testItem := range testData {
		t.Run(testItem.value, func(t *testing.T) {
			assert.Equal(t, testItem.expected, expandEnv(testItem.value))
		})
	}
}

func TestExpandEnvWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Test designed for Windows")
	}
	testData := []struct {
		value    string
		expected string
	}{
		{"%windir%", os.Getenv("windir")},
		{"%windir%%OS%", os.Getenv("windir") + os.Getenv("OS")},
	}

	for _, testItem := range testData {
		t.Run(testItem.value, func(t *testing.T) {
			assert.Equal(t, testItem.expected, expandEnv(testItem.value))
		})
	}
}
