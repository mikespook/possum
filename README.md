Possum
======

[![Build Status][travis-img]][travis]
[![GoDoc][godoc-img]][godoc]
[![Coverage Status](https://coveralls.io/repos/mikespook/possum/badge.svg?branch=master&service=github)](https://coveralls.io/github/mikespook/possum?branch=master)

Possum is a micro web library for Go.

It has following modules:

 * Routers
 * Views
 * Session
 * Helpers

Install
=======

Install the package:

```bash
go get github.com/mikespook/possum
```

Usage
=====

Importing the package and sub-packages:

```go
import (
	"github.com/mikespook/possum"
	"github.com/mikespook/possum/router"
	"github.com/mikespook/possum/view"
)
```

Possum uses `Context` for passing data, handling request and rendering response.

This is how to create a new server mux for Possum:

```go
mux := possum.NewServerMux()
```

And assign a customized error handler:

```go
mux.ErrorHandle = func(err error) {
	fmt.Println(err)
}
```

`PreRequest` and `PostResponse` are useful for pre-checking or customizing logs:

```go
mux.PreRequest = func(ctx *possum.Context) error {
	host, port, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		return err
	}
	if host != "127.0.0.1" {
		return possum.NewError(http.StatusForbidden, "Localhost only")
	}
	return nil
}

mux.PostResponse = func(ctx *possum.Context) error {
	fmt.Printf("[%d] %s:%s \"%s\"", ctx.Response.Status,
		ctx.Request.RemoteAddr,	ctx.Request.Method,
		ctx.Request.URL.String())		
}
```

A specific path can bind to a different combination of routers, handlers and views:

```go
f := session.NewFactory(session.CookieStorage('session-id', nil))

func helloword(ctx *Context) error {
	ctx.StartSession(f)
	return nil
}

mux.HandlerFunc(router.Simple("/json"), helloword, view.Json(view.CharSetUTF8))

if err := view.InitHtmlTemplates("*.html"); err != nil {
	return
}
mux.HandleFunc(router.Wildcard("/html/*/*"),
	helloworld, view.Html("base.html", "utf-8"))

if err := view.InitWatcher("*.html", view.InitTextTemplates, nil);
	err != nil {
	return
}
mux.HandleFunc(router.RegEx("/html/(.*)/[a-z]"),
	helloworld, view.Text("base.html", "utf-8"))

mux.HandleFunc(router.Colon("/:img/:id"), 
	nil, view.File("img.jpg", "image/jpeg"))
```

Also, a PProf methods can be initialized by `mux.InitPProf`:

```go
mux.InitPProf("/_pprof")
```

It will serve profiles and debug informations through `http://ip:port/_pprof`.

E.g.:

![][pprof]

And finally, it is a standard way for listening and serving:

```go
http.ListenAndServe(":8080", mux)
```

For more details, please see the [demo][demo].

Contributors
============

(_Alphabetic order_)
 
 * [Xing Xing][blog] <mikespook@gmail.com> [@Twitter][twitter]

Open Source - MIT Software License
==================================

See LICENSE.

 [travis-img]: https://travis-ci.org/mikespook/possum.png?branch=master
 [travis]: https://travis-ci.org/mikespook/possum
 [blog]: http://mikespook.com
 [twitter]: http://twitter.com/mikespook
 [godoc-img]: https://godoc.org/github.com/mikespook/gorbac?status.png
 [godoc]: https://godoc.org/github.com/mikespook/possum
 [coveralls-img]: https://coveralls.io/repos/mikespook/possum/badge.svg?branch=master&service=github
 [coveralls]: https://coveralls.io/github/mikespook/possum?branch=master
 [demo]: https://github.com/mikespook/possum/tree/master/demo
 [pprof]: https://pbs.twimg.com/media/CE4k3SIUMAAZiLy.png
