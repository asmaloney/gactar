package amod

import (
	"os"
)

func generateToStdout(str string) {
	_, log, _ := GenerateModel(str)
	log.Write(os.Stdout)
}

func Example_modelAuthors() {
	generateToStdout(`
	==model==
	name: Test
	authors {
    	'Andy Maloney <andy@example.com>' // did all the work
    	'Hiro Protagonist <hiro@example.com>' // fixed some things
	}
	==config==
	==init==
	==productions==`)

	// Output:
}

func Example_modelExamples() {
	generateToStdout(`
	==model==
	name: Test
	examples { [foo: bar] }
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==`)

	// Output:
}

func Example_modelExampleBadChunk() {
	generateToStdout(`
	==model==
	name: Test
	examples { [foo: bar] }
	==config==
	==init==
	==productions==`)

	// Output:
	// ERROR: could not find chunk named 'foo' (line 4)
}

func Example_gactarUnrecognizedField() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	gactar { foo: bar }
	==init==
	==productions==`)

	// Output:
	// ERROR: unrecognized field in gactar section: 'foo' (line 5)
}

func Example_chunkReservedName() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [_internal: foo bar] }
	==init==
	==productions==`)

	// Output:
	// ERROR: cannot use reserved chunk name '_internal' (chunks begining with '_' are reserved) (line 5)
}

func Example_chunkDuplicateName() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks {
    	[something: foo bar]
    	[something: foo bar]
	}
	==init==
	==productions==`)

	// Output:
	// ERROR: duplicate chunk name: 'something' (line 7)
}

func Example_modules() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	modules {
		imaginal { delay: 0.2 }
	}
	==init==
	==productions==`)

	// Output:
}

func Example_modulesUnrecognized() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	modules {
		foo { delay: 0.2 }
	}
	==init==
	==productions==`)

	// Output:
	// ERROR: unrecognized module in config: 'foo' (line 6)
}

func Example_imaginalFields() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	modules {
		imaginal { delay: 0.2 }
		memory { latency: 0.5 }
	}
	==init==
	==productions==`)

	// Output:
}

func Example_imaginalFieldType() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	modules {
		imaginal { delay: "gack" }
	}
	==init==
	==productions==`)

	// Output:
	// ERROR: imaginal delay 'gack' must be a number (line 6)
}

func Example_imaginalFieldRange() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	modules {
		imaginal { delay: -0.5 }
	}
	==init==
	==productions==`)

	// Output:
	// ERROR: imaginal delay '-0.500000' must be a positive number (line 6)
}

func Example_imaginalFieldUnrecognized() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	modules {
		imaginal { foo: bar }
	}
	==init==
	==productions==`)

	// Output:
	// ERROR: unrecognized field 'foo' in imaginal config (line 6)
}

func Example_memoryFieldUnrecognized() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	modules {
		memory { foo: bar }
	}
	==init==
	==productions==`)

	// Output:
	// ERROR: unrecognized field 'foo' in memory (line 6)
}

func Example_initializer1() {
	generateToStdout(`
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
	==productions==`)

	// Output:
}

func Example_initializer2() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [author: person object year] }
	==init==
	goal [author: Fred Book 1972]
	==productions==`)

	// Output:
}

func Example_initializer3() {
	// memory with one init is allowed
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [author: person object year] }
	==init==
	memory [author: Jane Book 1982]
	==productions==`)

	// Output:
}

func Example_initializerInvalidSlots() {
	// Check invalid number of slots in init
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [author: person object year] }
	==init==
	memory { [author: me software] }
	==productions==`)

	// Output:
	// ERROR: invalid chunk - 'author' expects 3 slots (line 7)
}

func Example_initializerInvalidChunk1() {
	// Check memory with invalid chunk
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	memory { [author: me software] }
	==productions==`)

	// Output:
	// ERROR: could not find chunk named 'author' (line 6)
}

func Example_initializerInvalidChunk2() {
	// Check buffer with invalid chunk
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	goal [author: Fred Book 1972]
	==productions==`)

	// Output:
	// ERROR: could not find chunk named 'author' (line 6)
}

func Example_initializerUnknownBuffer() {
	// Check unknown buffer
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [author: person object year] }
	==init==
	something [author: Fred Book 1972]
	==productions==`)

	// Output:
	// ERROR: buffer or memory 'something' not found in initialization  (line 7)
}

func Example_initializerMultipleInits() {
	// Check buffer with multiple inits
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [author: person object year] }
	==init==
	goal { [author: Fred Book 1972] [author: Jane Book 1982] }
	==productions==`)

	// Output:
	// ERROR: buffer 'goal' should only have one pattern in initialization (line 7)
}

func Example_productionUnusedVar() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: ?blat] }
		do { set goal to [foo: ding] }
	}`)

	// Output:
	// INFO: variable ?blat is not used - should be simplified to '?' (line 9)
}

func Example_productionInvalidAnonVarInSet1() {
	// https://github.com/asmaloney/gactar/issues/57
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: ?] }
		do { set goal.thing to ? }
	}`)

	// Output:
	// ERROR: cannot set 'goal.thing' to anonymous var ('?') in production 'start' (line 10)
}

func Example_productionInvalidAnonVarInSet2() {
	// https://github.com/asmaloney/gactar/issues/57
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: ?] }
		do { set goal to [foo: ?] }
	}`)

	// Output:
	// ERROR: cannot set 'goal.thing' to anonymous var ('?') in production 'start' (line 10)
}

func Example_productionUnusedVar2() {
	// Check that using a var twice in a buffer match does not get
	// marked as unused.
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match { goal [foo: ?blat ?blat] }
		do { set goal to [foo: ding] }
	}`)

	// Output:
}

func Example_productionInvalidMemory() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { another_goal [add: ? ?one1 ? ?one2 ? None?ans ?] }
		do { print 'foo' }
	}`)

	// Output:
	// ERROR: buffer 'another_goal' not found in production 'start' (line 8)
}

func Example_productionClearStatemtent() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { clear goal }
	}`)

	// Output:
}

func Example_productionClearStatemtentInvalidBuffer() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { clear some_buffer }
	}`)

	// Output:
	// ERROR: buffer 'some_buffer' not found in production 'start' (line 10)
}

func Example_productionSetStatemtentPattern() {
	// Check setting to pattern
	generateToStdout(`
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
	}`)

	// Output:
}

func Example_productionSetStatemtentVar() {
	// Check setting to var
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: ?blat] }
		do { set goal.thing to ?blat }
	}`)

	// Output:
}

func Example_productionSetStatemtentNil() {
	// Check setting to nil
	generateToStdout(`
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
			goal [foo: blat]
			imaginal [ack: bar]
		}
		do {
			set goal.thing to nil
			set imaginal.knowledge to nil
		}
	}`)

	// Output:
}

func Example_productionSetStatemtentNonVar() {
	// Check setting to non-existent var
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { set goal.thing to ?ding }
	}`)

	// Output:
	// ERROR: set statement variable '?ding' not found in matches for production 'start' (line 10)
}

func Example_productionSetStatemtentAssignNonPattern() {
	// Check setting buffer to non-pattern
	// https://github.com/asmaloney/gactar/issues/28
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { set goal to 6 }
	}`)

	// Output:
	// ERROR: buffer 'goal' must be set to a pattern in production 'start' (line 10)
}

func Example_productionSetStatemtentAssignNonsense() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { set goal to blat }
	}`)

	// Output:
	// ERROR: unexpected token "blat" (expected (SetValue | Pattern)) (line 10)
}

func Example_productionSetStatemtentAssignPattern() {
	// Check assignment of pattern to slot
	// https://github.com/asmaloney/gactar/issues/17
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { set goal.thing to [foo: ding] }
	}`)

	// Output:
	// ERROR: cannot set a slot ('thing') to a pattern in match buffer 'goal' in production 'start' (line 10)
}

func Example_productionRecallStatemtent() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match { goal [foo: ?next ?] }
		do { recall [foo: ?next ?] }
	}`)

	// Output:
}

func Example_productionRecallStatemtentVarNotFound() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { recall [count: ?next ?] }
	}`)

	// Output:
	// ERROR: recall statement variable '?next' not found in matches for production 'start' (line 10)
}

func Example_productionMultipleStatement() {
	generateToStdout(`
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
	}`)

	// Output:
}

func Example_productionChunkNotFound() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { goal [foo: error] }
		do { print 42 }
	}`)

	// Output:
	// ERROR: could not find chunk named 'foo' (line 8)
}

func Example_productionPrintStatement1() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match { goal [foo: ? ?other] }
		do { print 42, ?other, "blat" }
	}`)

	// Output:
}

func Example_productionPrintStatement2() {
	// Test print with vars from two different buffers
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match {
			goal [foo: ? ?other]
			retrieval [foo: ? ?other1]
		}
		do { print 42, ?other, ?other1, "blat" }
	}`)

	// Output:
}

func Example_productionPrintStatement3() {
	// print without args
	// https://github.com/asmaloney/gactar/issues/7
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: error] }
		do { print }
	}`)

	// Output:
}

func Example_productionPrintStatementInvalidID() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: error] }
		do { print fooID }
	}`)

	// Output:
	// ERROR: cannot use ID 'fooID' in print statement (line 9)
}

func Example_productionPrintStatementInvalidVar() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: error] }
		do { print ?fooVar }
	}`)

	// Output:
	// ERROR: print statement variable '?fooVar' not found in matches for production 'start' (line 9)
}

func Example_productionPrintStatementAnonymousVar() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: error] }
		do { print ? }
	}`)

	// Output:
	// ERROR: cannot print anonymous var ('?') in production 'start' (line 9)
}

func Example_productionMatchInternal() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: error] }
		do { print 42 }
	}`)

	// Output:
}

func Example_productionMatchInternalSlots() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: busy error] }
		do { print 42 }
	}`)

	// Output:
	// ERROR: _status should only have one slot for 'retrieval' in production 'start' (should be 'busy', 'empty', 'error', 'full') (line 8)
}

func Example_productionMatchInternalInvalidStatus1() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { goal [_status: something] }
		do { print 42 }
	}`)

	// Output:
	// ERROR: invalid _status 'something' for 'goal' in production 'start' (should be 'busy', 'empty', 'error', 'full') (line 8)
}

func Example_productionMatchInternalInvalidStatus2() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: something] }
		do { print 42 }
	}`)

	// Output:
	// ERROR: invalid _status 'something' for 'retrieval' in production 'start' (should be 'busy', 'empty', 'error', 'full') (line 8)
}
