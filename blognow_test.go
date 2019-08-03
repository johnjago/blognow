package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestParse(t *testing.T) {
	fileContents := `---
title = "My First Post"
date = 2019-07-20T00:00:00Z
---

# Hello
This is my first blog post.
`
	expected := Post{
		Title:   "My First Post",
		Date:    time.Date(2019, time.July, 20, 0, 0, 0, 0, time.UTC),
		Content: "<h3>Hello</h3>\n<p>This is my first blog post.</p>\n",
	}
	actual := parse(fileContents)
	if actual != expected {
		t.Errorf("parse(%q) returned %q, want %q", fileContents, actual, expected)
	}
}

func TestSlug(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"My First Post", "my-first-post"},
		{"Hello 世界", "hello-世界"},
		{"", ""},
	}
	for _, c := range cases {
		got := slug(c.in)
		if got != c.want {
			t.Errorf("slug(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCreateFile(t *testing.T) {
	name := "xidf9.test"
	contents := "This is the contents of the file."
	createFile(name, contents)
	readFileAndCompare(name, contents, t)
	os.Remove(name)
}

func readFileAndCompare(path, expected string, t *testing.T) {
	file, err := os.Open(path)
	check(err)
	data := make([]byte, 1000)
	bytesRead, err := file.Read(data)
	check(err)
	actual := string(data[:bytesRead])
	if actual != expected {
		t.Errorf(path + " does not contain the expected default contents.")
	}
}
