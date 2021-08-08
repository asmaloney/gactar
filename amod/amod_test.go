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

func TestProductionInvalidMemory(t *testing.T) {
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
		do {
			print 'foo'
		}
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

func TestProductionClearBuffer(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	buffers { bar }
	==productions==
	start {
		match {
			bar: 'foo'
		}
		do {
			clear some_buffer
		}
	}`

	_, err := GenerateModel(src)

	expected := "buffer not found in production 'start': 'some_buffer' (line 12)"
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
	buffers { bar, blat }
	==productions==
	start {
		match {
			bar: 'foo'
		}
		do {
			clear bar, blat
		}
	}`

	_, err = GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
