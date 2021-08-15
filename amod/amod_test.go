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
	checkExpectedError(err, expected, t)
}

func TestChunkReservedName(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks {
    	_internal( foo bar )
	}
	==productions==`

	_, err := GenerateModel(src)

	expected := "cannot use reserved chunk name '_internal' (chunks begining with '_' are reserved) (line 6)"
	checkExpectedError(err, expected, t)
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
	checkExpectedError(err, expected, t)
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
	checkExpectedError(err, expected, t)
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
	checkExpectedError(err, expected, t)

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
			set thing of goal to ` + "`foo ding`" + `
		}
	}`

	_, err = GenerateModel(src)

	expected := "cannot set a slot ('thing') to a pattern in match buffer 'goal' in production 'start' (line 14)"
	checkExpectedError(err, expected, t)
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
	checkExpectedError(err, expected, t)

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

func TestProductionChunkNotFound(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	==productions==
	start {
		match {
			goal: ` + "`foo error`" + `
		}
		do {
			print 42
		}
	}`

	_, err := GenerateModel(src)

	expected := "could not find chunk named 'foo' (line 8)"
	checkExpectedError(err, expected, t)
}

func TestProductionMatchInternal(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	==productions==
	start {
		match {
			memory: ` + "`_status error`" + `
		}
		do {
			print 42
		}
	}`

	_, err := GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	src = `
	==model==
	name: Test
	==config==
	==productions==
	start {
		match {
			memory: ` + "`_status busy error`" + `
		}
		do {
			print 42
		}
	}`

	_, err = GenerateModel(src)
	expected := "_status should only have one slot for 'memory' in production 'start' (should be 'busy', 'free', or 'error') (line 8)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	==productions==
	start {
		match {
			goal: ` + "`_status something`" + `
		}
		do {
			print 42
		}
	}`

	_, err = GenerateModel(src)
	expected = "invalid _status 'something' for 'goal' in production 'start' (should be 'full' or 'empty') (line 8)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	==productions==
	start {
		match {
			memory: ` + "`_status something`" + `
		}
		do {
			print 42
		}
	}`

	_, err = GenerateModel(src)
	expected = "invalid _status 'something' for 'memory' in production 'start' (should be 'busy', 'free', or 'error') (line 8)"
	checkExpectedError(err, expected, t)
}

func checkExpectedError(err error, expected string, t *testing.T) {
	t.Helper()

	if err == nil {
		t.Errorf("Expected error: %s", expected)
	} else {
		if err.Error() != expected {
			t.Errorf("Expected '%s' but got '%s'", expected, err.Error())
		}
	}
}
