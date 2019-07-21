# blognow ![](https://travis-ci.org/johnjago/blognow.svg?branch=master)
Blognow is a dead simple static site generator for blogs.

## Usage

`blognow my-blog` creates a new blog in the directory my-blog.

`blognow` generates a static site in dist/. You can copy these files to any
place where you can host static websites. It looks in the posts/ directory
and formats any correctly structured .md file as HTML.

As of now, it creates an index page (the latest post), individual post pages,
and an archive page.

The format for a post is:

```
---
title = "Post Title"
date = 2019-06-28
---

# Heading 1
## Heading 2
### Heading 3

This is a paragraph.

- This
- is
- a
- list.
```

Refer to a [Markdown reference](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet)
for a complete list of formatting options.

## Running your blog locally

```
$ npm i http-server -g
$ http-server dist/
```

## Contribute

Blognow is currently under development and is not stable. If you would like to
report an issue or suggest a feature, check out the GitHub issues.

Pull requests are also welcome!

## License
Unlicense
