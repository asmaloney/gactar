package errorlist

import (
	"fmt"
	"strings"
)

type Errors []string

// ErrorOrNil return the error if any, or nil if none.
// Used for return statements.
func (el Errors) ErrorOrNil() error {
	if len(el) == 0 {
		return nil
	}

	return el
}

// Add adds a string to the error list.
func (el *Errors) Add(e string) {
	*el = append(*el, e)
}

// Addf uses printf formatting to add an error.
func (el *Errors) Addf(e string, a ...interface{}) {
	el.Add(fmt.Sprintf(e, a...))
}

// String satifies the string interface to print errors joined by newlines.
func (el Errors) String() string {
	return strings.Join(el, "\n")
}

// Error satifies the error interface.
func (el Errors) Error() string {
	return el.String()
}
