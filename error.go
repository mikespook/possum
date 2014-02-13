package possum

import (
	"fmt"
)

type apiErr struct {
	status  int
	message string
}

func (err apiErr) Error() string {
	return err.message
}

func Errorf(status int, format string, a ...interface{}) apiErr {
	return apiErr{
		status:  status,
		message: fmt.Sprintf(format, a...),
	}
}
