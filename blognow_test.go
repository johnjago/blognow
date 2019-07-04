package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMakeBlogDir(t *testing.T) {
	basepath := "test-blog"
	paths := []string{
		filepath.Join(basepath, "posts"),
		filepath.Join(basepath, "posts/sample.md"),
		filepath.Join(basepath, "config.toml"),
	}
	makeBlogDir(basepath)

	// All sample files should exist.
	for _, path := range paths {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			t.Errorf(path + " should exist; it does not.")
		}
	}

	// Check the default files for the correct contents.
	readFileAndCompare(paths[1], samplePost, t)
	readFileAndCompare(paths[2], sampleConfig, t)

	os.RemoveAll(basepath)
}

func readFileAndCompare(path, expected string, t *testing.T) {
	file, err := os.Open(path)
	check(err)
	data := make([]byte, 100)
	bytesRead, err := file.Read(data)
	check(err)
	actual := string(data[:bytesRead])
	if actual != expected {
		t.Errorf(path + " does not contain the expected default contents.")
	}
}
