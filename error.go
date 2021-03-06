package possum

import (
	"errors"
	"log"
	"net/http"
)

var (
	errMethodNotAllowed = errors.New("possum: method not allowed")
)

func handleError(w http.ResponseWriter, errPanic error) {
	e, ok := errPanic.(Error)
	var status int
	var message string
	if ok {
		status = e.Status
	} else {
		status = http.StatusInternalServerError
	}
	message = e.Error()
	// use ErrorHandler to re-rander error output
	w.WriteHeader(status)
	if _, err := w.Write([]byte(message)); err != nil {
		log.Panicln(err)
	}
}

// Error implements error interface
type Error struct {
	Status  int
	Message string
}

func (err Error) Error() string {
	return err.Message
}
