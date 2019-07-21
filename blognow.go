/*
This is the CLI for Blognow.
*/
package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

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
const postsDir string = "posts"
const outputDir string = "dist/"

type BlogInfo struct {
	BaseURL string
	Title   string
	Tagline string
}

type Post struct {
	Title   string
	Date    time.Time
	Slug    string
	Content template.HTML
}

type PostPageData struct {
	Blog BlogInfo
	Post Post
}

type ArchivePageData struct {
	Blog  BlogInfo
	Posts []Post
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Building your blog...")
		build()
		fmt.Printf("Done. Output is in %s\n", outputDir)
		os.Exit(0)
	}
	makeBlogDir(os.Args[1])
}

// makeBlogDir creates the initial directory with sample files.
func makeBlogDir(path string) {
	postsPath := filepath.Join(path, postsDir)
	os.MkdirAll(postsPath, os.ModePerm)
	fmt.Println("Created a new blog: " + path)
	createFile(path+"/config.toml", sampleConfig)
	createFile(postsPath+"/sample.md", samplePost)
}

// build collects all the necessary information from configuration files
// and post files and builds the static site.
func build() {
	os.Mkdir(outputDir, os.ModePerm)

	// Get blog title, tagline, etc. from config.toml
	blogInfo := blogInfo()

	fmap := template.FuncMap{
		"formatDate": formatDate,
	}
	baseTemplate, err := template.New("").Funcs(fmap).ParseFiles(
		"templates/base.html",
		"templates/header.html",
	)
	check(err)
	postTemplate, err := template.Must(baseTemplate.Clone()).ParseFiles(
		"templates/post.html",
	)
	check(err)

	// Iterate over all .md files in posts/
	posts := make([]Post, 0)
	mostRecent := Post{}
	err = godirwalk.Walk(postsDir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			ext := filepath.Ext(osPathname)
			if ext == ".md" {
				// post.md -> Post struct
				content, err := ioutil.ReadFile(osPathname)
				check(err)
				post := parse(string(content))
				post.Slug = slug(post.Title)
				posts = append(posts, post)
				if post.Date.After(mostRecent.Date) {
					mostRecent = post
				}

				// Build output using template.
				data := PostPageData{
					Blog: blogInfo,
					Post: post,
				}
				var postHTML bytes.Buffer
				err = postTemplate.ExecuteTemplate(&postHTML, "base", data)
				check(err)
				os.Mkdir(outputDir+post.Slug, os.ModePerm)
				createFile(outputDir+post.Slug+"/index.html", postHTML.String())
			}
			return nil
		},
		Unsorted: true,
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", postsDir, err)
	}

	// Generate archive page.
	tmpl, err := template.Must(baseTemplate.Clone()).ParseFiles(
		"templates/archive.html",
	)
	check(err)
	archivePageData := ArchivePageData{
		Blog:  blogInfo,
		Posts: posts,
	}
	var archiveHTML bytes.Buffer
	err = tmpl.ExecuteTemplate(&archiveHTML, "base", archivePageData)
	check(err)
	os.Mkdir(outputDir+"archive", os.ModePerm)
	createFile(outputDir+"archive"+"/index.html", archiveHTML.String())

	// Use the most recent post as the index page.
	postPageData := PostPageData{
		Blog: blogInfo,
		Post: mostRecent,
	}
	var postHTML bytes.Buffer
	err = postTemplate.ExecuteTemplate(&postHTML, "base", postPageData)
	check(err)
	createFile(outputDir+"index.html", postHTML.String())
}

// blogInfo reads a config.toml file and returns a BlogInfo struct.
func blogInfo() BlogInfo {
	config, err := ioutil.ReadFile("config.toml")
	check(err)
	blogInfo := BlogInfo{}
	_, err = toml.Decode(string(config), &blogInfo)
	check(err)
	return blogInfo
}

// parse creates a Post struct from a post file containing front matter and
// Markdown.
func parse(content string) Post {
	post := Post{}
	frontMatter, err := extractFrontMatter(content)
	check(err)
	_, err = toml.Decode(frontMatter, &post)
	check(err)
	post.Content = template.HTML(extractBody(content))
	return post
}

// extractFrontMatter returns the front matter as a string given
// the entire contents of a post file.
func extractFrontMatter(content string) (string, error) {
	frontMatter := ""
	lines := strings.Split(content, "\n")
	if len(lines) < 2 {
		return frontMatter, errors.New("Error: Post file missing front matter")
	}
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			break
		} // End of front matter
		frontMatter += lines[i] + "\n"
	}
	return frontMatter, nil
}

// extractBody takes the content of a post file, converts the Markdown in the
// content to an HTML string, and returns this string.
func extractBody(content string) string {
	body := ""
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if i < 4 {
			continue
		} // Skip past the front matter
		// Since the blog and post title (h1 and h2, respectively) are added
		// automatically, all other headings start two levels down.
		if strings.HasPrefix(line, "#") {
			line = "##" + line
		}
		body += line + "\n"
	}
	md := markdown.New()
	return md.RenderToString([]byte(body))
}

func createFile(name string, content string) {
	contentBytes := []byte(content)
	err := ioutil.WriteFile(name, contentBytes, 0644)
	check(err)
}

// slug turns a string into this-kind-of-format that can be used in a URL.
func slug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
