/*
This is the CLI for Blognow.
*/
package main

import (
	"bytes"
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

type BlogInfo struct {
	BaseURL string
	Title   string
	Tagline string
}

type Post struct {
	Title   string
	Date    time.Time
	Slug	string
	Content template.HTML
}

type PostPageData struct {
	Blog BlogInfo
	Post Post
}

type AllPostsPageData struct {
	Blog BlogInfo
	Posts []Post
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
	postsPath := filepath.Join(path, postsDir)
	os.MkdirAll(postsPath, os.ModePerm)
	fmt.Println("Created a new blog: " + path)
	createFile(path+"/config.toml", sampleConfig)
	createFile(postsPath+"/sample.md", samplePost)
}

func build() {
	outputDir := "dist/"
	os.Mkdir(outputDir, os.ModePerm)

	// Get blog title, tagline, etc. from config.toml
	blogInfo := blogInfo();

	// Iterate over all .md files in posts/
	posts := make([]Post, 0)
	err := godirwalk.Walk(postsDir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			ext := filepath.Ext(osPathname)
			if ext == ".md" {
				// post.md -> Post struct
				content, err := ioutil.ReadFile(osPathname)
				check(err)
				post := parse(string(content))
				post.Slug = slug(post.Title)
				posts = append(posts, post)

				// Build out the template
				tmpl := template.Must(template.ParseFiles(
					"templates/base.html",
					"templates/post.html",
					"templates/header.html",
				))
				templateData := PostPageData{
					Blog: blogInfo,
					Post: post,
				}
				var singlePostHTML bytes.Buffer
				err = tmpl.ExecuteTemplate(&singlePostHTML, "base", templateData)
				check(err)

				// Generate output file
				os.Mkdir(outputDir+post.Slug, os.ModePerm)
				createFile(outputDir+post.Slug+"/index.html", singlePostHTML.String())
			}
			return nil
		},
		Unsorted: true,
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", postsDir, err)
	}

	// Generate all posts page
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/archive.html",
		"templates/header.html",
	))
	templateData := AllPostsPageData{
		Blog: blogInfo,
		Posts: posts,
	}
	var allPostsHTML bytes.Buffer
	err = tmpl.ExecuteTemplate(&allPostsHTML, "base", templateData)
	check(err)
	os.Mkdir(outputDir+"archive", os.ModePerm)
	createFile(outputDir+"archive"+"/index.html", allPostsHTML.String())

	// get the most recent post and copy it to dist/index.html
}

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
		// TODO: error handle
		return ""
	}
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			break
		} // End of front matter
		frontMatter += lines[i] + "\n"
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
