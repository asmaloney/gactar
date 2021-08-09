package amod

import (
	"strings"
	"testing"
)

func TestMinimumModel(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	==productions==`

	_, err := parse(strings.NewReader(src))

	if err != nil {
		t.Errorf("Could not parse minimal src: %s", err.Error())
	}
}
