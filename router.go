package possum

import (
	"container/list"
	"net/url"
	"sync"

	"github.com/mikespook/possum/router"
	"github.com/mikespook/possum/view"
)

// Routers contains all routers.
type Routers struct {
	sync.RWMutex
	s map[string]struct {
		r router.Router
		v view.View
		h HandlerFunc
	}
	l *list.List
}

// Find a router with the specific path and return it.
func (rs *Routers) Find(path string) (url.Values, HandlerFunc, view.View) {
	defer rs.RUnlock()
	rs.RLock()
	if s, ok := rs.s[path]; ok {
		return nil, s.h, s.v
	}
	for e := rs.l.Front(); e != nil; e = e.Next() {
		s := e.Value.(struct {
			r router.Router
			v view.View
			h HandlerFunc
		})
		if params, ok := s.r.Match(path); ok {
			return params, s.h, s.v
		}
	}
	return nil, nil, nil
}

// Add a router to list
func (rs *Routers) Add(r router.Router, h HandlerFunc, v view.View) {
	defer rs.Unlock()
	rs.Lock()
	s := struct {
		r router.Router
		v view.View
		h HandlerFunc
	}{r, v, h}
	// simple will full-match the path
	if sr, ok := r.(*router.Base); ok {
		rs.s[sr.Path] = s
		return
	}
	rs.l.PushFront(s)
}

// NewRouters initailizes Routers instance.
func NewRouters() *Routers {
	return &Routers{
		s: make(map[string]struct {
			r router.Router
			v view.View
			h HandlerFunc
		}),
		l: list.New(),
	}
}
