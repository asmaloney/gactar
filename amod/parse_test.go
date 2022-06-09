package amod

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alecthomas/participle/v2"
)

func TestMinimumModel(t *testing.T) {
	t.Parallel()

	src := `
	==model==
	name: Test
	==config==
	==init==
	==productions==`

	_, err := parse(strings.NewReader(src))

	if err != nil {
		t.Errorf("Could not parse minimal src: %s", err.Error())
	}
}

func FuzzExampleModels(f *testing.F) {
	match, err := filepath.Glob("../examples/*.amod")
	if err != nil {
		f.Fatal(err)
	}

	for _, input := range match {
		code, err := os.ReadFile(input)
		if err != nil {
			f.Fatal(err)
		}

		f.Add(string(code))
	}

	f.Fuzz(func(t *testing.T, orig string) {
		_, err := parse(strings.NewReader(orig))
		if err != nil &&
			!strings.Contains(err.Error(), "must match at least once") && // participle.parseError is not public, so hack it for now...
			!errors.As(err, &participle.UnexpectedTokenError{}) &&
			!errors.As(err, &LexError{}) {
			t.Errorf("Error: %s\nInput: %q", err.Error(), orig)
		}
	})
}
