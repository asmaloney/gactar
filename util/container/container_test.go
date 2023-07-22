package container

import "testing"

func TestGetIndex1(t *testing.T) {
	t.Parallel()

	list := []string{"a", "b", "c", "d"}

	index := GetIndex1("a", list)
	if index != 1 {
		t.Errorf("Incorrect index: expected 1")
	}

	// Check case-sensitivity
	index = GetIndex1("C", list)
	if index == 3 {
		t.Errorf("Incorrect index: expected -1")
	}

	index = GetIndex1("d", list)
	if index != 4 {
		t.Errorf("Incorrect index: expected 4")
	}

	index = GetIndex1("x", list)
	if index != -1 {
		t.Errorf("Incorrect return: expected -1")
	}

}

func TestUniqueAndSorted(t *testing.T) {
	t.Parallel()

	list := []string{"d", "b", "c", "b", "a", "c"}

	list2 := UniqueAndSorted(list)

	if len(list2) != 4 {
		t.Errorf("Items not removed from list")
	}

	expected := []string{"a", "b", "c", "d"}
	for i, v := range expected {
		if v != list2[i] {
			t.Errorf("Resulting list incorrect")
		}
	}
}

func TestFindAndDelete(t *testing.T) {
	t.Parallel()

	list := []string{"a", "b", "c", "b", "d", "c"}

	list2 := FindAndDelete(list, "b")

	if len(list2) != 4 {
		t.Errorf("Items not removed from list")
	}

	expected := []string{"a", "c", "d", "c"}
	for i, v := range expected {
		if v != list2[i] {
			t.Errorf("Resulting list incorrect")
		}
	}
}
