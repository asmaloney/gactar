package errorlist

import (
	"fmt"
	"strings"
)

type Errors []string

func (el *Errors) Add(e string) {
	*el = append(*el, e)
}

func (el *Errors) Addf(e string, a ...interface{}) {
	el.Add(fmt.Sprintf(e, a...))
}

func (el Errors) String() string {
	return strings.Join(el, "\n")
}

func (el Errors) Error() string {
	return el.String()
}
