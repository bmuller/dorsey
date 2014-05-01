# Dorsey: A Go Microframework

[![GoDoc](https://godoc.org/github.com/bmuller/dorsey?status.png)](https://godoc.org/github.com/bmuller/dorsey)
[![Build Status](https://travis-ci.org/bmuller/dorsey.png?branch=master)](https://travis-ci.org/bmuller/dorsey)

Dorsey is a micro framework for [Go](http://golang.org).  Did the world need another?  Probably not, but Dorsey does handle filtering, context, and authentication in a slightly more streamlined way than others (though much of the additional delight is syntax more than anything else).

## Installation

```bash
go get github.com/bmuller/dorsey
```

## Examples

Dorsey is simple.

```go
func showContent(w *dorsey.ResponseWriter, r *dorsey.Request) {
     w.Render("Hi There " + w.GetURLParam("name"))
}

server := dorsey.New()
server.Get("/users/:name", showContent)
server.Run(":8080")
```

Dorsey can render lots of stuff, like files.

```go
func showAFile(w *dorsey.ResponseWriter, r *dorsey.Request) {
     w.RenderFile("/tmp/a/file.txt")
}

server := dorsey.New()
server.Get("/afile", showAFile)
server.Run(":8080")
```

Dorsey has special handlers to help you write less code.

```go
server := dorsey.New()

// Render a file
server.Get("/afile", dorsey.FileHandler("/tmp/a/file.txt")

// Render a directory
server.Get("/static/", dorsey.DirectoryHandler("/tmp/static")

// Redirect somewhere
server.Get("/oldpath", dorsey.RedirectHandler("/newpath"))

server.Run(":8080")
```

Dorsey doesn't need before filters - each handler can act as a filter naturally.  Attach many to a path, and as soon as one renders or redirects then the rest won't be called.

```go
func checkUser(w *dorsey.ResponseWriter, r *dorsey.Request) {
     if r.GetParam("username", "") != "validuser" {
          w.Render("You're not authorized!")
     }
}

func secretContent(w *dorsey.ResponseWriter, r *dorsey.Request) {
     // only shown to valid users
     w.Render("A SECRET")
}

server.Get("/secret", checkUser, secretContent)
```

Dorsey can even handle HTTP basic auth for you with a special handler.

```go
func checkUser(username, password string) {
     return username == password
}

// This basicAuth handler can be used for all authentication
basicAuth := dorsey.BasicAuthHandler(auth)

server := dorsey.New()
server.Get("/secret/one", basicAuth, secretOne)
server.Get("/secret/two", basicAuth, secretTwo)
```

Finally, Dorsey can handle passing context information from one handler to the next.

```go
func setUser(w *dorsey.ResponseWriter, r *dorsey.Request) {
     w.Context["username"] = r.GetParam("username", "")
}

func showContent(w *dorsey.ResponseWriter, r *dorsey.Request) {
     w.Render("Hi There " + w.Context["username"])
}

server.Get("/secret", checkUser, secretContent)
```
