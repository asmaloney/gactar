package amod

func Example_productionClearStatement() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'blat'] }
		do { clear goal }
	}`)

	// Output:
}

func Example_productionErrorClearStatementInvalidBuffer() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'blat'] }
		do { clear some_buffer }
	}`)

	// Output:
	// ERROR: buffer 'some_buffer' not found in production 'start' (line 10, col 7)
}

func Example_productionSetStatementPattern() {
	// Check setting to pattern
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		description: "This is a description"
		match { goal [foo: 'blat'] }
		do { set goal to [foo: 'ding'] }
	}`)

	// Output:
}

func Example_productionSetStatementVar() {
	// Check setting to var
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?blat] }
		do { set goal.thing to ?blat }
	}`)

	// Output:
}

func Example_productionSetStatementID() {
	// Check setting to ID
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'thing'] }
		do { set goal.thing to thing2 }
	}`)

	// Output:
}

func Example_productionSetStatementString() {
	// Check setting to ID
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'thing'] }
		do { set goal.thing to 'thing string' }
	}`)

	// Output:
}

func Example_productionSetStatementNil() {
	// Check setting to nil
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		imaginal { delay: 0.2 }
	}
	chunks {
		[foo: thing]
		[ack: knowledge]
	}
	~~ init ~~
	~~ productions ~~
	start {
		match {
			goal [foo: 'blat']
			imaginal [ack: 'bar']
		}
		do {
			set goal.thing to nil
			set imaginal.knowledge to nil
		}
	}`)

	// Output:
}

func Example_productionErrorSetStatementNonBuffer() {
	// Check setting to non-existent buffer in set statement
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'blat'] }
		do { set foo.bar to 'blat' }
	}`)

	// Output:
	// ERROR: buffer "foo" not found in model (line 10, col 11)
}

func Example_productionErrorSetStatementNonBuffer2() {
	// Check setting to buffer not used in the match
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules { imaginal { delay: 0.2 } }
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'blat'] }
		do { set imaginal.bar to 'blat' }
	}`)

	// Output:
	// ERROR: match buffer 'imaginal' not found in production 'start' (line 11, col 11)
}

func Example_productionErrorSetStatementInvalidSlot() {
	// Check setting to buffer not used in the match
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'blat'] }
		do { set goal.bar to 'blat' }
	}`)

	// Output:
	// ERROR: slot 'bar' does not exist in chunk type 'foo' for match buffer 'goal' in production 'start' (line 10, col 16)
}

func Example_productionErrorSetStatementNonVar1() {
	// Check setting to non-existent var
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'blat'] }
		do { set goal.thing to ?ding }
	}`)

	// Output:
	// ERROR: set statement variable '?ding' not found in matches for production 'start' (line 10, col 25)
}

func Example_productionErrorSetStatementNonVar2() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'blat'] }
		do { set goal to [foo: ?ding] }
	}`)

	// Output:
	// ERROR: set statement variable '?ding' not found in matches for production 'start' (line 10, col 25)
}

func Example_productionErrorSetStatementAssignNonPattern() {
	// Check setting buffer to non-pattern
	// https://github.com/asmaloney/gactar/issues/28
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'blat'] }
		do { set goal to 6 }
	}`)

	// Output:
	// ERROR: buffer 'goal' must be set to a pattern in production 'start' (line 10, col 19)
}

func Example_productionErrorSetStatementAssignPattern() {
	// Check assignment of pattern to slot
	// https://github.com/asmaloney/gactar/issues/17
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'blat'] }
		do { set goal.thing to [foo: 'ding'] }
	}`)

	// Output:
	// ERROR: cannot set a slot ('goal.thing') to a pattern in production 'start' (line 10, col 11)
}

func Example_productionRecallStatement() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next *] }
	}`)

	// Output:
}

func Example_productionErrorRecallStatementMultiple() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next *] }
		do {
			recall [foo: ?next *]
			recall [foo: * ?next]
		}
	}`)

	// Output:
	// ERROR: only one recall statement per production is allowed in production 'start' (line 12, col 3)
}

func Example_productionErrorRecallStatementInvalidPattern() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next * 'bar'] }
	}`)

	// Output:
	// ERROR: invalid chunk - 'foo' expects 2 slots (line 10, col 14)
}

func Example_productionErrorRecallStatementVarNotFound() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] [bar: other thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: 'blat'] }
		do { recall [bar: ?next *] }
	}`)

	// Output:
	// ERROR: recall statement variable '?next' not found in matches for production 'start' (line 10, col 20)
}

func Example_productionRecallStatementWithWith() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next *] with ( recently_retrieved reset ) }
	}`)

	// Output:
}

func Example_productionRecallStatementWithMultipleWith() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next *] with ( recently_retrieved reset ) and ( recently_retrieved t ) }
	}`)

	// Output:
}

func Example_productionErrorRecallStatementWithInvalidParam() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next *] with ( foo_param 42 ) }
	}`)

	// Output:
	// ERROR: recall 'with': unrecognized option "foo_param". (line 10, col 34)
}

func Example_productionRecallStatementWithNIL() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next *] with ( recently_retrieved nil ) }
	}`)

	// Output:
}

func Example_productionRecallStatementWithID() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next *] with ( recently_retrieved t ) }
	}`)

	// Output:
}

func Example_productionErrorRecallStatementWithString() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next *] with ( recently_retrieved "bar" ) }
	}`)

	// Output:
	// ERROR: recall 'with': invalid value "bar" for option "recently_retrieved" (expected one of: t, nil, reset). (line 10, col 34)
}

func Example_productionErrorRecallStatementWithVar() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next *] with ( recently_retrieved ?bar ) }
	}`)

	// Output:
	// ERROR: recall 'with': parameter 'recently_retrieved'. Unexpected variable (line 10, col 34)
}

func Example_productionErrorRecallStatementWithIncorrectNumArgs() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next *] with ( recently_retrieved nil nil ) }
	}`)

	// Output:
	// ERROR: unexpected token "nil" (expected ")") (line 10, col 59)
}
func Example_productionMultipleStatement() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?next ?other] }
		do {
        	recall [foo: ?next *]
			set goal to [foo: ?other 42]
		}
	}`)

	// Output:
}

func Example_productionErrorChunkNotFound() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: error] }
		do { print 42 }
	}`)

	// Output:
	// ERROR: could not find chunk named 'foo' (line 8, col 16)
}

func Example_productionPrintStatement1() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: * ?other] }
		do { print 42, ?other, 'blat' }
	}`)

	// Output:
}

func Example_productionPrintStatement2() {
	// Test print with vars from two different buffers
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match {
			goal [foo: * ?other]
			retrieval [foo: * ?other1]
		}
		do { print 42, ?other, ?other1, 'blat' }
	}`)

	// Output:
}

func Example_productionPrintStatement3() {
	// print without args
	// https://github.com/asmaloney/gactar/issues/7
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { buffer_state goal empty }
		do { print }
	}`)

	// Output:
}

func Example_productionPrintStatementBuffer() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { buffer_state retrieval empty }
		do { print goal }
	}`)

	// Output:
}

func Example_productionPrintStatementBufferSlot() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: * *] }
		do { print goal.thing1 }
	}`)

	// Output:
}

func Example_productionErrorPrintStatementInvalidBuffer() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { buffer_state retrieval empty }
		do { print fooID }
	}`)

	// Output:
	// ERROR: buffer "fooID" not found in model (line 9, col 13)
}

func Example_productionErrorPrintStatementInvalidBufferSlot() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: * *] }
		do { print goal.blat }
	}`)

	// Output:
	// ERROR: slot 'blat' does not exist in chunk type 'foo' for match buffer 'goal' in production 'start' (line 10, col 18)
}

func Example_productionErrorPrintStatementInvalidNil() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { buffer_state goal empty }
		do { print nil }
	}`)

	// Output:
	// ERROR: unexpected token "nil" (expected "}") (line 9, col 13)
}

func Example_productionErrorPrintStatementInvalidVar() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { buffer_state retrieval empty }
		do { print ?fooVar }
	}`)

	// Output:
	// ERROR: print statement variable '?fooVar' not found in matches for production 'start' (line 9, col 13)
}
