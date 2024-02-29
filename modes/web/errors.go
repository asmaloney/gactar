package web

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyRequestBody = errors.New("empty request body")
	ErrNoModel          = errors.New("no model loaded")
)

type ErrInvalidModelID struct {
	ID int
}

func (e ErrInvalidModelID) Error() string {
	return fmt.Sprintf("invalid model id: %d", e.ID)
}

type ErrInvalidSessionID struct {
	ID int
}

func (e ErrInvalidSessionID) Error() string {
	return fmt.Sprintf("invalid session id: %d", e.ID)
}
