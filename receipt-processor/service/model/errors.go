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
	return fmt.Errorf("error %d - %s", NotFound, cause)
}

func ErrBadRequest(cause string) error {
	return fmt.Errorf("error %d - %s", BadRequest, cause)
}

func ErrInternalServer(cause string) error {
	return fmt.Errorf("unexpected error %d - %s", Internal, cause)
}
