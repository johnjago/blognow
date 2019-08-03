package main

const baseTmpl string = `{{define "base"}}
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/kognise/water.css@latest/dist/light.min.css">
  <style>
	.center {
		text-align: center;
	}
  </style>
  <title>{{template "title" .}}</title>
</head>
<body>
  {{template "header" .}}
  {{template "body" .}}
</body>
</html>
{{end}}
`

const headerTmpl string = `{{define "header"}}
  <div class="center">
	<h1><a href="/">{{.Blog.Title}}</a></h1>
	<p>{{.Blog.Tagline}}</p>
  </div>
{{end}}
`

const postTmpl string = `{{define "title"}}
  {{.Post.Title}} - {{.Blog.Title}}
{{end}}

{{define "body"}}
  <h2>{{.Post.Title}}</h2>
  {{.Post.Content}}
  <hr>
  <p>Posted: {{.Post.Date | formatDate}}</p>
  <a href="/archive/">See all posts »</a>
{{end}}
`

const archiveTmpl string = `{{define "title"}}
All Posts - {{.Blog.Title}}
{{end}}

{{define "body"}}
<ul>
  {{range .Posts}}
  <li>{{.Date | formatDate}} <a href="/{{.Slug}}">{{.Title}}</a></li>
  {{end}}
</ul>
{{end}}
`
