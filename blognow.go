package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const sampleConfig string = "baseURL = \"https://example.org/\"\ntitle = \"My Blog\"\n"
const samplePost string = "---\ntitle: \"My First Post\"\ndate: 2019-05-05\n---\n"

func main() {
	if len(os.Args) == 1 {
		// TODO: Build site
		fmt.Println("Building your blog")
		os.Exit(0)
	}

	// Create a new blog
	path := os.Args[1]
	postsPath := filepath.Join(path, "posts")
	os.MkdirAll(postsPath, os.ModePerm)
	fmt.Println("Created blog " + path)
	createFile(path+"/config.toml", sampleConfig)
	createFile(postsPath+"/sample.md", samplePost)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createFile(name string, content string) {
	contentBytes := []byte(content)
	err := ioutil.WriteFile(name, contentBytes, 0644)
	check(err)
}
