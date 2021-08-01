package amod

import (
	"testing"
)

func TestACTRUnrecognizedField(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	actr { foo: bar }
	==init==
	==productions==`

	_, err := GenerateModel(src)

	expected := "Unrecognized field in actr section: 'foo'"
	if err == nil {
		t.Errorf("Expected error: %s", expected)
	} else {
		if err.Error() != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, err.Error())
		}
	}
}
