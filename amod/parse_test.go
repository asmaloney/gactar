package amod

import (
	"strings"
	"testing"
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
