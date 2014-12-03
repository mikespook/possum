package possum

import (
	"container/list"
	"regexp"
	"strings"
	"sync"
)

// Router is an interface to match specific path.
type Router interface {
	Match(string) bool
	View() View
	Handler() HandlerFunc
	HandleFunc(HandlerFunc, View)
}

// Routers contains all routers.
type Routers struct {
	sync.RWMutex
	s map[string]Router
	l *list.List
}

// Find a router with the specific path and return it.
func (r *Routers) Find(path string) Router {
	defer r.RUnlock()
	r.RLock()
	if router, ok := r.s[path]; ok {
		return router
	}
	for e := r.l.Front(); e != nil; e = e.Next() {
		router := e.Value.(Router)
		if router.Match(path) {
			return router
		}
	}
	return nil
}

// Add a router to list
func (r *Routers) Add(router Router) {
	defer r.Unlock()
	r.Lock()

	// SimpleRouter will full-match the path
	if simpleRouter, ok := router.(*SimpleRouter); ok {
		r.s[simpleRouter.path] = router
		return
	}
	r.l.PushFront(router)
}

func NewRouters() *Routers {
	return &Routers{
		s: make(map[string]Router),
		l: list.New(),
	}
}

type SimpleRouter struct {
	path    string
	view    View
	handler HandlerFunc
}

func NewSimpleRouter(path string) *SimpleRouter {
	return &SimpleRouter{
		path: path,
	}
}

func (r *SimpleRouter) Match(path string) bool {
	return path == r.path
}

func (r *SimpleRouter) View() View {
	return r.view
}

func (r *SimpleRouter) Handler() HandlerFunc {
	return r.handler
}

func (r *SimpleRouter) HandleFunc(handler HandlerFunc, view View) {
	r.handler = handler
	r.view = view
}

// WildcardRouter matches path with wildcard
type WildcardRouter struct {
	SimpleRouter
	matches []string
}

func NewWildcardRouter(path string) *WildcardRouter {
	matches := strings.Split(path, "/")
	return &WildcardRouter{
		SimpleRouter: SimpleRouter{
			path: path,
		},
		matches: matches,
	}
}

func (r *WildcardRouter) Match(path string) bool {
	matches := strings.Split(path, "/")
	if len(matches) != len(r.matches) {
		return false
	}
	for k, v := range r.matches {
		if v != "*" && matches[k] != v {
			return false
		}
	}
	return true
}

// RegExRouter matches path with regex
type RegExRouter struct {
	SimpleRouter
	r *regexp.Regexp
}

func NewRegExRouter(path string) *RegExRouter {
	r, err := regexp.Compile(path)
	if err != nil {
		panic(err)
	}
	return &RegExRouter{
		SimpleRouter: SimpleRouter{
			path: path,
		},
		r: r,
	}
}

func (r *RegExRouter) Match(path string) bool {
	return r.r.MatchString(path)
}

// ResourceRouter matches path with REST-full resources form.
type ResourceRouter struct {
	SimpleRouter
	matches []string
}

func NewResourceRouter(path string) *ResourceRouter {
	matches := strings.Split(path, "/")
	return &ResourceRouter{
		SimpleRouter: SimpleRouter{
			path: path,
		},
		matches: matches,
	}
}

func (r *ResourceRouter) Match(path string) bool {
	matches := strings.Split(path, "/")
	i := 0
	for _, v := range r.matches {
		if v != "" && v[0] == ':' {
			if matches[i] == v[1:] {
				i++
			} else {
				return false
			}
		} else if matches[i] != v {
			return false
		}
		i++
	}
	return true
}
