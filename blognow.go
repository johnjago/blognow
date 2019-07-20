/*
This is the CLI for Blognow.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"html/template"
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/karrick/godirwalk"
	"gitlab.com/golang-commonmark/markdown"
)

const sampleConfig string = `
baseURL = "https://example.org/"
title = "My Blog"
tagline = "Don't sail too close to the wind"
`

const samplePost string = `
---
title = "My First Post"
date = 2019-05-05
---
`

type Post struct {
	Title   string
	Date    time.Time
	Content template.HTML
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Building your blog...")
		build()
		fmt.Println("Done. Output is in dist/")
		os.Exit(0)
	}
	makeBlogDir(os.Args[1])
}

func makeBlogDir(path string) {
	postsPath := filepath.Join(path, "posts")
	os.MkdirAll(postsPath, os.ModePerm)
	fmt.Println("Created a new blog: " + path)
	createFile(path+"/config.toml", sampleConfig)
	createFile(postsPath+"/sample.md", samplePost)
}

func createFile(name string, content string) {
	contentBytes := []byte(content)
	err := ioutil.WriteFile(name, contentBytes, 0644)
	check(err)
}

func build() {
	postsDir := "posts"
	err := godirwalk.Walk(postsDir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			ext := filepath.Ext(osPathname)
			if ext == ".md" {
				filename := filepath.Base(strings.TrimSuffix(osPathname, ext))
				fmt.Println(filename)
				content, err := ioutil.ReadFile(osPathname)
				check(err)
				post := parse(string(content))
				fmt.Println(post)
				tmpl := template.Must(template.ParseFiles("templates/post.html"))
				var buff bytes.Buffer
				tmpl.Execute(&buff, post)
				fmt.Println(buff.String())
				// output template to file
			}
			return nil
		},
		Unsorted: true,
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", postsDir, err)
	}

	// for each post in posts/
	// parse content into a data structure
	// insert into a template
	// output the template as new html file in dist/post-title/index.html
	// get the most recent post and copy it to dist/index.html
	// iterate through post list and create all posts page
}

// parse creates a Post struct from a post file containing front matter and
// Markdown.
func parse(content string) Post {
	post := Post{}
	frontMatter := extractFrontMatter(content)
	_, err := toml.Decode(frontMatter, &post)
	check(err)
	post.Content = template.HTML(extractBody(content))
	return post
}

// extractFrontMatter returns the front matter as a string given
// the entire contents of a post file.
func extractFrontMatter(content string) string {
	frontMatter := ""
	lines := strings.Split(content, "\n")
	if len(lines) < 2 {
		// error handle
		return ""
	}
	for i, line := range lines {
		if i == 0 {
			continue
		}
		if line == "---" {
			break
		}
		frontMatter += line + "\n"
	}
	return frontMatter
}

// extractBody takes the content of a post file, converts the Markdown in the
// content to an HTML string, and returns this string.
func extractBody(content string) string {
	body := ""
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if i < 4 {
			continue
		}
		body += line + "\n"
	}
	md := markdown.New()
	return md.RenderToString([]byte(body))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
