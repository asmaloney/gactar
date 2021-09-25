package amod

import (
	"testing"
)

func TestModelExamples(t *testing.T) {
	src := `
	==model==
	name: Test
	examples { [foo: bar] }
	==config==
	==init==
	==productions==`

	_, err := GenerateModel(src)

	expected := "could not find chunk named 'foo' (line 4)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	examples { [foo: bar] }
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==`

	_, err = GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestGACATRUnrecognizedField(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	gactar { foo: bar }
	==init==
	==productions==`

	_, err := GenerateModel(src)

	expected := "unrecognized field in gactar section: 'foo' (line 5)"
	checkExpectedError(err, expected, t)
}

func TestChunkReservedName(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks { [_internal: foo bar] }
	==init==
	==productions==`

	_, err := GenerateModel(src)

	expected := "cannot use reserved chunk name '_internal' (chunks begining with '_' are reserved) (line 5)"
	checkExpectedError(err, expected, t)
}

func TestChunkDuplicateName(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks {
    	[something: foo bar]
    	[something: foo bar]
	}
	==init==
	==productions==`

	_, err := GenerateModel(src)

	expected := "duplicate chunk name: 'something' (line 7)"
	checkExpectedError(err, expected, t)
}

func TestModules(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	modules {
		imaginal { delay: 0.2 }
	}
	==init==
	==productions==`

	_, err := GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	src = `
	==model==
	name: Test
	==config==
	modules {
		foo { delay: 0.2 }
	}
	==init==
	==productions==`

	_, err = GenerateModel(src)

	expected := "unrecognized module in config: 'foo' (line 6)"
	checkExpectedError(err, expected, t)
}

func TestImaginalFields(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	modules {
		imaginal { delay: 0.2 }
		memory { latency: 0.5 }
	}
	==init==
	==productions==`

	_, err := GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	src = `
	==model==
	name: Test
	==config==
	modules {
		imaginal { delay: "gack" }
	}
	==init==
	==productions==`

	_, err = GenerateModel(src)

	expected := "imaginal delay 'gack' must be a number (line 6)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	modules {
		imaginal { delay: -0.5 }
	}
	==init==
	==productions==`

	_, err = GenerateModel(src)

	expected = "imaginal delay '-0.500000' must be a positive number (line 6)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	modules {
		imaginal { foo: bar }
	}
	==init==
	==productions==`

	_, err = GenerateModel(src)

	expected = "unrecognized field 'foo' in imaginal config (line 6)"
	checkExpectedError(err, expected, t)
}

func TestMemoryUnrecognizedField(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	modules {
		memory { foo: bar }
	}
	==init==
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
		[remember: person]
		[author: person object]
	}
	==init==
	memory {
		[remember: me]
		[author: me software]
	}
	==productions==`

	_, err := GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	src = `
	==model==
	name: Test
	==config==
	chunks { [author: person object year] }
	==init==
	goal [author: Fred Book 1972]
	==productions==`

	_, err = GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Check invalid number of slots in init
	src = `
	==model==
	name: Test
	==config==
	chunks { [author: person object year] }
	==init==
	memory { [author: me software] }
	==productions==`

	_, err = GenerateModel(src)

	expected := "invalid chunk - 'author' expects 3 slots (line 7)"
	checkExpectedError(err, expected, t)

	// Check memory with invalid chunk
	src = `
	==model==
	name: Test
	==config==
	==init==
	memory { [author: me software] }
	==productions==`

	_, err = GenerateModel(src)

	expected = "could not find chunk named 'author' (line 6)"
	checkExpectedError(err, expected, t)

	// Check buffer with invalid chunk
	src = `
	==model==
	name: Test
	==config==
	==init==
	goal [author: Fred Book 1972]
	==productions==`

	_, err = GenerateModel(src)
	checkExpectedError(err, expected, t)

	// Check unknown buffer
	src = `
	==model==
	name: Test
	==config==
	chunks { [author: person object year] }
	==init==
	something [author: Fred Book 1972]
	==productions==`

	_, err = GenerateModel(src)
	expected = "buffer or memory 'something' not found in initialization  (line 7)"
	checkExpectedError(err, expected, t)

	// Check buffer with multiple inits
	src = `
	==model==
	name: Test
	==config==
	chunks { [author: person object year] }
	==init==
	goal { [author: Fred Book 1972] [author: Jane Book 1982] }
	==productions==`

	_, err = GenerateModel(src)
	expected = "buffer 'goal' should only have one pattern in initialization (line 7)"
	checkExpectedError(err, expected, t)

	// memory with one init is allowed
	src = `
	==model==
	name: Test
	==config==
	chunks { [author: person object year] }
	==init==
	memory [author: Jane Book 1982]
	==productions==`

	_, err = GenerateModel(src)
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
		match { another_goal [add: ? ?one1 ? ?one2 ? None?ans ?] }
		do { print 'foo' }
	}`

	_, err := GenerateModel(src)

	expected := "buffer 'another_goal' not found in production 'start' (line 8)"
	checkExpectedError(err, expected, t)
}

func TestProductionClearStatement(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { clear some_buffer }
	}`

	_, err := GenerateModel(src)

	expected := "buffer 'some_buffer' not found in production 'start' (line 10)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { clear goal }
	}`

	_, err = GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestProductionSetStatement(t *testing.T) {
	// Check setting to pattern
	src := `
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		description: "This is a description"
		match { goal [foo: blat] }
		do { set goal to [foo: ding] }
	}`

	_, err := GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Check setting to var
	src = `
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: ?blat] }
		do { set goal.thing to ?blat }
	}`

	_, err = GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Check setting to nil
	src = `
	==model==
	name: Test
	==config==
	modules {
		imaginal { delay: 0.2 }
	}
	chunks {
		[foo: thing]
		[ack: knowledge]
	}
	==init==
	==productions==
	start {
		match {
			goal [foo: ?blat]
			imaginal [ack: ?bar]
		}
		do {
			set goal.thing to nil
			set imaginal.knowledge to nil
		}
	}`

	_, err = GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Check setting to non-existent var
	src = `
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: ?blat] }
		do { set goal.thing to ?ding }
	}`

	_, err = GenerateModel(src)
	expected := "set statement variable '?ding' not found in matches for production 'start' (line 10)"
	checkExpectedError(err, expected, t)

	// https://github.com/asmaloney/gactar/issues/28
	src = `
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { set goal to 6 }
	}`

	_, err = GenerateModel(src)

	expected = "buffer 'goal' must be set to a pattern in production 'start' (line 10)"
	checkExpectedError(err, expected, t)

	// https://github.com/asmaloney/gactar/issues/17
	src = `
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { set goal.thing to [foo: ding] }
	}`

	_, err = GenerateModel(src)

	expected = "cannot set a slot ('thing') to a pattern in match buffer 'goal' in production 'start' (line 10)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { set goal to blat }
	}`

	_, err = GenerateModel(src)

	expected = `10:22: unexpected token "blat" (expected (SetValue | Pattern))`
	checkExpectedError(err, expected, t)
}

func TestProductionRecallStatement(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { recall [count: ?next ?] }
	}`

	_, err := GenerateModel(src)

	expected := "recall statement variable '?next' not found in matches for production 'start' (line 10)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match { goal [foo: ?next ?other] }
		do { recall [foo: ?next ?] }
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
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match { goal [foo: ?next ?other] }
		do {
        	recall [foo: ?next ?]
			set goal to [foo: ?other 42]
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
	==init==
	==productions==
	start {
		match { goal [foo: error] }
		do { print 42 }
	}`

	_, err := GenerateModel(src)

	expected := "could not find chunk named 'foo' (line 8)"
	checkExpectedError(err, expected, t)
}

func TestProductionPrintStatement(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match { goal [foo: ?next ?other] }
		do { print 42, ?other, "blat" }
	}`

	_, err := GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Test print with vars from two different buffers
	src = `
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match {
			goal [foo: ?next ?other]
			retrieval [foo: ?next1 ?other1]
		}
		do { print 42, ?other, ?other1, "blat" }
	}`

	_, err = GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// https://github.com/asmaloney/gactar/issues/7
	src = `
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: error] }
		do { print }
	}`

	_, err = GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	src = `
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: error] }
		do { print fooID }
	}`

	_, err = GenerateModel(src)
	expected := "cannot use ID 'fooID' in print statement (line 9)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: error] }
		do { print ?fooVar }
	}`

	_, err = GenerateModel(src)
	expected = "print statement variable '?fooVar' not found in matches for production 'start' (line 9)"
	checkExpectedError(err, expected, t)
}

func TestProductionMatchInternal(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: error] }
		do { print 42 }
	}`

	_, err := GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	src = `
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: busy error] }
		do { print 42 }
	}`

	_, err = GenerateModel(src)
	expected := "_status should only have one slot for 'retrieval' in production 'start' (should be 'busy', 'empty', 'error', 'full') (line 8)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { goal [_status: something] }
		do { print 42 }
	}`

	_, err = GenerateModel(src)
	expected = "invalid _status 'something' for 'goal' in production 'start' (should be 'busy', 'empty', 'error', 'full') (line 8)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: something] }
		do { print 42 }
	}`

	_, err = GenerateModel(src)
	expected = "invalid _status 'something' for 'retrieval' in production 'start' (should be 'busy', 'empty', 'error', 'full') (line 8)"
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
