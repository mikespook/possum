Possum
======

[![Build Status][travis-img]][travis]

Possum is a micro web library for Go.

_Possum has not been officially released yet, as it is still in active development._

Install
=======

Install the package:

```bash
go get github.com/mikespook/possum
```

Usage
=====

Possum uses `Context` passing data, handling request and rendering response.

Create a new possum server mux :

```go
	mux := possum.NewServerMux()
```

Assign a customized error handler:

```go
mux.ErrorHandler = func(err error) {
	fmt.Println(err)
}
```

`PreRequest` and `PostResponse` is useful for pre-checking or customizing logs:

```go
mux.PreRequest = func(ctx *Context) error {
	host, port, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		return err
	}
	if host != "127.0.0.1" {
		return possum.NewError(http.StatusForbidden, "Localhost only")
	}
	return nil
}

mux.PostResponse = func(ctx *Context) error {
	fmt.Printf("[%d] %s:%s \"%s\"", ctx.Response.Status, ctx.Request.RemoteAddr,
		ctx.Request.Method, ctx.Request.URL.String())		
}
```

Add handlers with different views:

```go
f := session.NewFactory(session.CookieStorage('session-id', nil))

func helloword(ctx *Context) error {
	ctx.StartSession(f)
	return nil
}

mux.HandlerFunc("/json", helloword, possum.JsonView{})

if err := possum.InitHtmlTemplates("*.html"); err != nil {
	return
}
mux.HandleFunc("/html", helloworld, possum.NewHtmlView("base.html", "utf-8"))

if err := possum.InitViewWatcher("*.html", possum.InitTextTemplates, nil); err != nil {
	return
}
mux.HandleFunc("/html", helloworld, possum.NewTextView("base.html", "utf-8"))

mux.HandleFunc("/img.jpg", nil, possum.NewFileView("img.jpg", "image/jpeg"))
```

Also, PProf can be initialized by `mux.InitPProf`:

```go
mux.InitPProf("/_pprof")
```

And finally, listen and serve:

```go
http.ListenAndServe(":8080", mux)
```

For more details, see [Demo][demo].

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
 [demo]: https://github.com/mikespook/possum/blob/master/demo/main.go
