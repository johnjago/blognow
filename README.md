# blognow
Blognow is a dead simple static site generator for blogs.

## Usage

`blognow my-blog` creates a new blog in the directory my-blog.

`blognow` generates a static site in dist/. You can copy these files to any
place where you can host static websites. It looks in the posts/ directory
and formats any correctly structured .md file as HTML.

The format for a post is:

```
---
title: "Post Title"
date: 2019-06-28
---

# Heading 1
## Heading 2
### Heading 3

This is a paragraph.
```

Refer to a [Markdown reference](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet)
for a complete list of formatting options.

## License
Unlicense
