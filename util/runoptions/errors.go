package runoptions

import "fmt"

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
