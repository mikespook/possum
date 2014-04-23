Possum
======

[![Build Status][travis-img]][travis]

Possum is a micro web-api framework for Go.

_Possum has not been officially released yet, as it is still in active development._

Install
=======

Install the package:

```bash
go get github.com/mikespook/possum
```

Usage
=====

Possum framework supply two kinds api methods:

 * RPC: A function binding to a path can be called through HTTP.
 * REST: A Resource(struct) binding to a path can be called through HTTP with deferent request methods for deferent usages.

The Possum's Handler implementing http.HandlerFunc interface can be set to http.ListenAndServe or http.ListenAndServeTLS.

Define a rpc function:

```go
// foobar responses intpu params.
func foobar(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
	return http.StatusOK, params
}
```

Define a resource implementing interfaces:

```go
type Foobar struct {
	data string
	possum.NoDelete
	possum.NoPatch
	possum.NoPost
}

func (foobar *Foobar) Get(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
	return http.StatusOK, foobar.data
}

func (foobar *Foobar) Put(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
	foobar.data = params.Get("data")
	return http.StatusOK, ""
}
```

Get a new handler of possum:

```go
	h := possum.NewHandler()
```

Assign a custome error handler:

```go
h.ErrorHandler = func(err error) {
	fmt.Println(err)
}
```

A wrap handler is usually used for global pre-checking or custome logs:

```go
h.PreHandler = func(r *http.Request) (int, error) {
	host, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if host != "127.0.0.1" {
		return http.StatusForbidden, fmt.Errorf("Localhost only")
	}
	return http.StatusOK, nil
}

h.PostHandler = func(r *http.Request, status int) {
	fmt.Printf("[%d] %s:%s \"%s\"", status, r.RemoteAddr, r.Method, r.URL.String())		
}
```

Bind the rpc function to a path:

```go
h.AddRPC("/rpc/test", foobar)
```

Bind the resource to a path:

```go
if err := h.AddResource("/rest/test", &Foobar{}); err != nil {
	fmt.Println(err)
	return
}
```

Listen and serve it:

```go
http.ListenAndServe(":8080", h)
```

You can add some wrap functions to a rpc directly:

```go
func checkSecret(handler possum.HandlerFunc) possum.HandlerFunc {
	return func(params url.Values) (status int, data interface{}) {
		if params.Get("secret") != secret {
			return http.StatusForbidden, fmt.Errorf("Wrong secret")
		}
		return handler(params)
	}
}

h.AddRPC("/rpc/test", checkSecret(a))
```

Resources need `Wrap` function:

```go
wrap, err := possum.Wrap(checkSecret, &Foobar{})
if err != nil {
	fmt.Println(err)
	return
}
if err := handler.AddResource("/rest/user", wrap); err != nil {
	fmt.Println(err)
	return
}
```

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
