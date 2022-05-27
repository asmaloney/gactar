package issues

import "testing"

func TestAddEntryLocationNil(t *testing.T) {
	log := New()

	log.Warning(nil, "test warning")

	i := log.issues[0]

	if i.Location != nil {
		t.Errorf("Expected location to be nil")
	}
}

func TestAddEntryLocationEmpty(t *testing.T) {
	log := New()

	log.Warning(&Location{}, "test warning")

	i := log.issues[0]

	if i.Location != nil {
		t.Errorf("Expected location to be nil")
	}
}
