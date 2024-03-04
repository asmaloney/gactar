package runoptions

import (
	"fmt"
	"strings"
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
	return fmt.Sprintf("invalid framework name: %q; expected one of %q", e.Name, strings.Join(ValidFrameworks, ", "))
}

type ErrInvalidLogLevel struct {
	Level string
}

func (e ErrInvalidLogLevel) Error() string {
	return fmt.Sprintf("invalid log level: %q; expected one of %q", e.Level, strings.Join(ACTRLoggingLevels, ", "))
}
