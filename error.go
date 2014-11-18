package possum

import "fmt"

// NewError returns an new Error given
// a response status code and a error msg.
func NewError(status int, msg string) Error {
	return Error{status, msg}
}

// An Error represents an error response to be send to client.
type Error struct {
	Status  int
	Message string
}

func (err Error) Error() string {
	return fmt.Sprintf("%s(%d)", err.Message, err.Status)
}
