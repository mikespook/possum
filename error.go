package possum

import "fmt"

type Error struct {
	Status  int
	Message string
}

func NewError(status int, msg string) Error {
	return Error{status, msg}
}

func (err Error) Error() string {
	return fmt.Sprintf("%s(%d)", err.Message, err.Status)
}
