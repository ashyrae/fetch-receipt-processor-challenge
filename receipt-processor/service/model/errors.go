package model

import (
	"fmt"
)

const (
	BadRequest = 400
	NotFound   = 404
	Internal   = 503
)

func ErrNotFound(cause string) error {
	return fmt.Errorf("error %d - Resource was not found: %s", NotFound, cause)
}

func ErrBadRequest(cause string) error {
	return fmt.Errorf("error %d - Request was invalid: %s", NotFound, cause)
}

func ErrInternalServer(cause string) error {
	return fmt.Errorf("error %d - Unexpected internal server error: %s", Internal, cause)
}
