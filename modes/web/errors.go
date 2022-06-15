package web

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyRequestBody = errors.New("empty request body")
	ErrNoModel          = errors.New("no model loaded")
)

type ErrFrameworkNotActive struct {
	Name string
}

func (e ErrFrameworkNotActive) Error() string {
	return fmt.Sprintf("framework %q is not active on server", e.Name)
}

type ErrInvalidFrameworkName struct {
	Name string
}

func (e ErrInvalidFrameworkName) Error() string {
	return fmt.Sprintf("invalid framework name: %q", e.Name)
}

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
