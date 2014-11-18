package session

import (
	"net/http"
)

type FactoryFunc func(http.ResponseWriter, *http.Request) (*Session, error)

func NewFactory(storage Storage) FactoryFunc {
	return func(w http.ResponseWriter, r *http.Request) (s *Session, err error) {
		s = &Session{
			storage: storage,
			w:       w,
		}
		if err = s.Init(); err != nil {
			return
		}
		if err = storage.LoadTo(r, s); err != http.ErrNoCookie {
			return
		}
		return s, nil
	}
}
