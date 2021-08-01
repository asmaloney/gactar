package errorlist

import "testing"

func TestErrors(t *testing.T) {
	var err Errors = nil

	if err.ErrorOrNil() != nil {
		t.Errorf("Expected nil: (%+v)", err)
	}

	expected := "some error"
	err.Add(expected)

	if err.ErrorOrNil() == nil {
		t.Errorf("Expected '%s', got nil", expected)
	}

	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}
