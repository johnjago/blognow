package main

const baseTmpl string = `{{define "base"}}
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/water.css@2/out/water.min.css">
  <style>
    .center {
      text-align: center;
    }
    .no-bullet {
      list-style-type: none;
    }
    .p0 {
      padding: 0;
    }
    .mr2 {
      margin-right: 2rem;
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
  <p>{{.Post.Date | formatDate}}</p>
  <a href="/archive/">See all posts »</a>
{{end}}
`

const archiveTmpl string = `{{define "title"}}
All Posts - {{.Blog.Title}}
{{end}}

{{define "body"}}
  {{range $year := .Years}}
    <h2>{{$year}}</h2>
      <ul class="p0">
        {{$posts := index $.YearGroups $year}}
        {{range $posts}}
          <li class="no-bullet">
            <span class="mr2">{{.Date | formatArchiveDate}}</span> <a href="/{{.Slug}}">{{.Title}}</a>
          </li>
        {{end}}
      </ul>
  {{end}}
{{end}}
`
