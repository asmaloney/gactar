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

	expected := "unrecognized field in actr section: 'foo' (line 5)"
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
	memory {
    	foo: bar
	}
	==productions==`

	_, err := GenerateModel(src)

	expected := "unrecognized field 'foo' in memory (line 6)"
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
	chunks {
		remember( who )
	}
	==init==
	memory {
		'remember me'
	}
	==productions==`

	_, err := GenerateModel(src)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestProductionInvalidMemory(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match {
			another_goal: ` + "`add ? ?one1 ? ?one2 ? None?ans ?`" + `
		}
		do {
			print 'foo'
		}
	}`

	_, err := GenerateModel(src)

	expected := "buffer or memory 'another_goal' not found in production 'start' (line 9)"
	if err == nil {
		t.Errorf("Expected error: %s", expected)
	} else {
		if err.Error() != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, err.Error())
		}
	}
}

func TestProductionClearStatement(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks {
		foo( thing )
	}
	==productions==
	start {
		match {
			goal: ` + "`foo blat`" + `
		}
		do {
			clear some_buffer
		}
	}`

	_, err := GenerateModel(src)

	expected := "buffer 'some_buffer' not found in production 'start' (line 14)"
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
	chunks {
		foo( thing )
	}
	==productions==
	start {
		match {
			goal: ` + "`foo blat`" + `
		}
		do {
			clear goal
		}
	}`

	_, err = GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestProductionSetStatement(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks {
		foo( thing )
	}
	==productions==
	start {
		match {
			goal: ` + "`foo blat`" + `
		}
		do {
			set goal to ` + "`foo ding`" + `
		}
	}`

	_, err := GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestProductionRecallStatement(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks {
		foo( thing )
	}
	==productions==
	start {
		match {
			goal: ` + "`foo blat`" + `
		}
		do {
        	recall ` + "`count ?next ?`" + `
		}
	}`

	_, err := GenerateModel(src)

	expected := "recall statement variable '?next' not found in matches for production 'start' (line 14)"
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
	chunks {
		foo( thing1 thing2 )
	}
	==productions==
	start {
		match {
			goal: ` + "`foo ?next ?other`" + `
		}
		do {
        	recall ` + "`foo ?next ?`" + `
		}
	}`

	_, err = GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestProductionMultipleStatement(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks {
		foo( thing1 thing2 )
	}
	==productions==
	start {
		match {
			goal: ` + "`foo ?next ?other`" + `
		}
		do {
        	recall ` + "`foo ?next ?`" + `
			set goal to ` + "`foo ?other 42`" + `
		}
	}`

	_, err := GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
