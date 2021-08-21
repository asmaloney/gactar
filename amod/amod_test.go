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
	chunks { _internal( foo bar ) }
	==productions==`

	_, err := GenerateModel(src)

	expected := "cannot use reserved chunk name '_internal' (chunks begining with '_' are reserved) (line 5)"
	checkExpectedError(err, expected, t)

	// https://github.com/asmaloney/gactar/issues/23
	src = `
	==model==
	name: Test
	==config==
	chunks { goal( foo bar ) }
	==productions==`

	_, err = GenerateModel(src)

	expected = "cannot use reserved chunk name 'goal' (line 5)"
	checkExpectedError(err, expected, t)
}

func TestChunkDuplicateName(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks {
    	something( foo bar )
    	something( foo bar )
	}
	==productions==`

	_, err := GenerateModel(src)

	expected := "duplicate chunk name: 'something' (line 7)"
	checkExpectedError(err, expected, t)
}

func TestMemoryUnrecognizedField(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	memory { foo: bar }
	==productions==`

	_, err := GenerateModel(src)

	expected := "unrecognized field 'foo' in memory (line 5)"
	checkExpectedError(err, expected, t)
}

func TestInitializers(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks {
		remember( person )
		author( person object )
	}
	==init==
	memory {
		` + "`remember(me)`" + `
		` + "`author(  me  software  )`" + `
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
	chunks { author( person object year ) }
	==init==
	memory { ` + "`author( me software )`" + ` }
	==productions==`

	_, err = GenerateModel(src)

	expected := "invalid initialization - expected 3 slots (line 7)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	==init==
	memory { ` + "`author( me software )`" + ` }
	==productions==`

	_, err = GenerateModel(src)

	expected = "could not find chunk named 'author' in initialization (line 6)"
	checkExpectedError(err, expected, t)
}

func TestProductionInvalidMemory(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { another_goal ` + "`add( ? ?one1 ? ?one2 ? None?ans ? )`" + ` }
		do { print 'foo' }
	}`

	_, err := GenerateModel(src)

	expected := "buffer or memory 'another_goal' not found in production 'start' (line 8)"
	checkExpectedError(err, expected, t)
}

func TestProductionClearStatement(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks { foo( thing ) }
	==productions==
	start {
		match { goal ` + "`foo( blat )`" + ` }
		do { clear some_buffer }
	}`

	_, err := GenerateModel(src)

	expected := "buffer 'some_buffer' not found in production 'start' (line 9)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	chunks { foo( thing ) }
	==productions==
	start {
		match { goal ` + "`foo( blat )`" + ` }
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
	chunks { foo( thing ) }
	==productions==
	start {
		match { goal ` + "`foo( blat )`" + ` }
		do { set goal to ` + "`foo( ding )`" + ` }
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
	chunks { foo( thing ) }
	==productions==
	start {
		match { goal ` + "`foo( ?blat )`" + ` }
		do { set thing of goal to ?blat }
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
	chunks { foo( thing ) }
	==productions==
	start {
		match { goal ` + "`foo( ?blat )`" + ` }
		do { set thing of goal to ?ding }
	}`

	_, err = GenerateModel(src)
	expected := "set statement variable '?ding' not found in matches for production 'start' (line 9)"
	checkExpectedError(err, expected, t)

	// https://github.com/asmaloney/gactar/issues/28
	src = `
	==model==
	name: Test
	==config==
	chunks { foo( thing ) }
	==productions==
	start {
		match { goal ` + "`foo( blat )`" + ` }
		do { set goal to 6 }
	}`

	_, err = GenerateModel(src)

	expected = "buffer 'goal' must be set to a pattern in production 'start' (line 9)"
	checkExpectedError(err, expected, t)

	// https://github.com/asmaloney/gactar/issues/17
	src = `
	==model==
	name: Test
	==config==
	chunks { foo( thing ) }
	==productions==
	start {
		match { goal ` + "`foo( blat )`" + ` }
		do { set thing of goal to ` + "`foo( ding )`" + ` }
	}`

	_, err = GenerateModel(src)

	expected = "cannot set a slot ('thing') to a pattern in match buffer 'goal' in production 'start' (line 9)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	chunks { foo( thing ) }
	==productions==
	start {
		match { goal ` + "`foo( blat )`" + ` }
		do { set goal to blat }
	}`

	_, err = GenerateModel(src)

	expected = `9:22: unexpected token "blat" (expected (SetValue | Pattern))`
	checkExpectedError(err, expected, t)
}

func TestProductionRecallStatement(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks { foo( thing ) }
	==productions==
	start {
		match { goal ` + "`foo( blat )`" + ` }
		do { recall ` + "`count( ?next ? )`" + ` }
	}`

	_, err := GenerateModel(src)

	expected := "recall statement variable '?next' not found in matches for production 'start' (line 9)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	chunks { foo( thing1 thing2 ) }
	==productions==
	start {
		match { goal ` + "`foo( ?next ?other )`" + ` }
		do { recall ` + "`foo( ?next ? )`" + ` }
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
	chunks { foo( thing1 thing2 ) }
	==productions==
	start {
		match { goal ` + "`foo( ?next ?other )`" + ` }
		do {
        	recall ` + "`foo( ?next ? )`" + `
			set goal to ` + "`foo( ?other 42 )`" + `
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
		match { goal ` + "`foo( error )`" + ` }
		do { print 42 }
	}`

	_, err := GenerateModel(src)

	expected := "could not find chunk named 'foo' (line 7)"
	checkExpectedError(err, expected, t)
}

func TestProductionPrintStatement(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	chunks { foo( thing1 thing2 ) }
	==productions==
	start {
		match { goal ` + "`foo( ?next ?other )`" + ` }
		do { print 42, ?other, "blat" }
	}`

	_, err := GenerateModel(src)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// https://github.com/asmaloney/gactar/issues/7
	src = `
	==model==
	name: Test
	==config==
	==productions==
	start {
		match { memory ` + "`_status( error )`" + ` }
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
	==productions==
	start {
		match { memory ` + "`_status( error )`" + ` }
		do { print fooID }
	}`

	_, err = GenerateModel(src)
	expected := "cannot use ID 'fooID' in print statement (line 8)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	==productions==
	start {
		match { memory ` + "`_status( error )`" + ` }
		do { print ?fooVar }
	}`

	_, err = GenerateModel(src)
	expected = "print statement variable '?fooVar' not found in matches for production 'start' (line 8)"
	checkExpectedError(err, expected, t)
}

func TestProductionMatchInternal(t *testing.T) {
	src := `
	==model==
	name: Test
	==config==
	==productions==
	start {
		match { memory ` + "`_status( error )`" + ` }
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
	==productions==
	start {
		match { memory ` + "`_status( busy error )`" + ` }
		do { print 42 }
	}`

	_, err = GenerateModel(src)
	expected := "_status should only have one slot for 'memory' in production 'start' (should be 'busy', 'free', or 'error') (line 7)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	==productions==
	start {
		match { goal ` + "`_status( something )`" + ` }
		do { print 42 }
	}`

	_, err = GenerateModel(src)
	expected = "invalid _status 'something' for 'goal' in production 'start' (should be 'full' or 'empty') (line 7)"
	checkExpectedError(err, expected, t)

	src = `
	==model==
	name: Test
	==config==
	==productions==
	start {
		match { memory ` + "`_status( something )`" + ` }
		do { print 42 }
	}`

	_, err = GenerateModel(src)
	expected = "invalid _status 'something' for 'memory' in production 'start' (should be 'busy', 'free', or 'error') (line 7)"
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
