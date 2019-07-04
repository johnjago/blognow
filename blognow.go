/*
This is the main CLI program for Blognow.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const sampleConfig string = `baseURL = "https://example.org/"
title = "My Blog"
`

const samplePost string = `---
title: "My First Post"
date: 2019-05-05
---
`

func main() {
	if len(os.Args) == 1 {
		// TODO: Build site
		fmt.Println("Building your blog")
		os.Exit(0)
	}
	makeBlogDir(os.Args[1])
}

func makeBlogDir(path string) {
	postsPath := filepath.Join(path, "posts")
	os.MkdirAll(postsPath, os.ModePerm)
	fmt.Println("Created blog " + path)
	createFile(path+"/config.toml", sampleConfig)
	createFile(postsPath+"/sample.md", samplePost)
}

func createFile(name string, content string) {
	contentBytes := []byte(content)
	err := ioutil.WriteFile(name, contentBytes, 0644)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
