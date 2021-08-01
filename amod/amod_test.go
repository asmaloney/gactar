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

	expected := "Unrecognized field in actr section: 'foo' (line 5)"
	if err == nil {
		t.Errorf("Expected error: %s", expected)
	} else {
		if err.Error() != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, err.Error())
		}
	}
}

func TestMemoryBufferField(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	memories {
    	a_memory { buffer: foo }
	}
	==init==
	==productions==`

	_, err := GenerateModel(src)

	expected := "buffer not found for memory 'a_memory': foo (line 6)"
	if err == nil {
		t.Errorf("Expected error: %s", expected)
	} else {
		if err.Error() != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, err.Error())
		}
	}

	src = `
	==model==
	name: Test
	==config==
	memories {
    	a_memory { buffer: 42 }
	}
	==init==
	==productions==`

	_, err = GenerateModel(src)

	expected = "buffer should not be a number in memory 'a_memory': 42 (line 6)"
	if err == nil {
		t.Errorf("Expected error: %s", expected)
	} else {
		if err.Error() != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, err.Error())
		}
	}
}

func TestMemoryUnrecognizedField(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	memories {
    	a_memory { foo: bar }
	}
	==init==
	==productions==`

	_, err := GenerateModel(src)

	expected := "Unrecognized field in memory 'a_memory': 'foo' (line 6)"
	if err == nil {
		t.Errorf("Expected error: %s", expected)
	} else {
		if err.Error() != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, err.Error())
		}
	}
}

func TestInitializers(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	buffers { bar }
	memories {
    	a_memory { buffer: bar }
	}
	==init==
	another_memory {
		'remember me'
	}
	==productions==`

	_, err := GenerateModel(src)

	expected := "memory not found for initialization 'another_memory' (line 10)"
	if err == nil {
		t.Errorf("Expected error: %s", expected)
	} else {
		if err.Error() != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, err.Error())
		}
	}
}

func TestProductions(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	buffers { bar }
	memories {
    	a_memory { buffer: bar }
	}
	==init==
	a_memory {
		'remember me'
	}
	==productions==
	start {
		match {
			another_goal: 'add ? ?one1 ? ?one2 ? None?ans ?'
		}
		do #<
			print('foo')
		>#
	}`

	_, err := GenerateModel(src)

	expected := "buffer or memory not found for production 'start': another_goal (line 16)"
	if err == nil {
		t.Errorf("Expected error: %s", expected)
	} else {
		if err.Error() != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, err.Error())
		}
	}
}
