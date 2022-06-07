package amod

func Example_production() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: ?blat] }
		do {
			print ?blat
			stop
		}
	}`)

	// Output:
}

func Example_productionWildcard() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match { goal [foo: ?blat *] }
		do { print ?blat }
	}`)

	// Output:
}

func Example_productionNotWildcard() {
	// Odd error message.
	// See: https://github.com/asmaloney/gactar/issues/124
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match { goal [foo: ?blat !*] }
		do { print ?blat }
	}`)

	// Output:
	// ERROR: unexpected token "!" (expected "]") (line 9, col 27)
}

func Example_productionUnusedVar1() {
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
	// ERROR: variable ?blat is not used - should be simplified to '*' (line 9, col 21)
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
		match { another_goal [add: * ?one1 * ?one2 * None?ans *] }
		do { print 'foo' }
	}`)

	// Output:
	// ERROR: buffer 'another_goal' not found in production 'start' (line 8, col 10)
}

func Example_productionClearStatement() {
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

func Example_productionClearStatementInvalidBuffer() {
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
	// ERROR: buffer 'some_buffer' not found in production 'start' (line 10, col 7)
}

func Example_productionSetStatementPattern() {
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

func Example_productionSetStatementVar() {
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

func Example_productionSetStatementNil() {
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

func Example_productionSetStatementNonBuffer() {
	// Check setting to non-existent buffer in set statement
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { set foo.bar to 'blat' }
	}`)

	// Output:
	// ERROR: buffer 'foo' not found (line 10, col 11)
	// ERROR: match buffer 'foo' not found in production 'start' (line 10, col 11)
}

func Example_productionSetStatementNonBuffer2() {
	// Check setting to buffer not used in the match
	generateToStdout(`
	==model==
	name: Test
	==config==
	modules { imaginal { delay: 0.2 } }
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { set imaginal.bar to 'blat' }
	}`)

	// Output:
	// ERROR: match buffer 'imaginal' not found in production 'start' (line 11, col 11)
}

func Example_productionSetStatementInvalidSlot() {
	// Check setting to buffer not used in the match
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { set goal.bar to 'blat' }
	}`)

	// Output:
	// ERROR: slot 'bar' does not exist in chunk 'foo' for match buffer 'goal' in production 'start' (line 10, col 16)
}

func Example_productionSetStatementNonVar1() {
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
	// ERROR: set statement variable '?ding' not found in matches for production 'start' (line 10, col 25)
}

func Example_productionSetStatementNonVar2() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { set goal to [foo: ?ding] }
	}`)

	// Output:
	// ERROR: set statement variable '?ding' not found in matches for production 'start' (line 10, col 25)
}

func Example_productionSetStatementCompoundVar() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] }
	==init==
	==productions==
	start {
		match { goal [foo: ?ding] }
		do { set goal to [foo: ?ding!5] }
	}`)

	// Output:
	// ERROR: cannot set 'goal.thing' to compound var in production 'start' (line 10, col 25)
}

func Example_productionSetStatementAssignNonPattern() {
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
	// ERROR: buffer 'goal' must be set to a pattern in production 'start' (line 10, col 19)
}

func Example_productionSetStatementAssignNonsense() {
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
	// ERROR: unexpected token "blat" (expected (SetValue | Pattern)) (line 10, col 19)
}

func Example_productionSetStatementAssignPattern() {
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
	// ERROR: cannot set a slot ('goal.thing') to a pattern in production 'start' (line 10, col 11)
}

func Example_productionRecallStatement() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next *] }
	}`)

	// Output:
}

func Example_productionRecallStatementMultiple() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
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

func Example_productionRecallStatementInvalidPattern() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing1 thing2] }
	==init==
	==productions==
	start {
		match { goal [foo: ?next *] }
		do { recall [foo: ?next * bar] }
	}`)

	// Output:
	// ERROR: invalid chunk - 'foo' expects 2 slots (line 10, col 14)
}

func Example_productionRecallStatementVarNotFound() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	chunks { [foo: thing] [bar: other thing] }
	==init==
	==productions==
	start {
		match { goal [foo: blat] }
		do { recall [bar: ?next *] }
	}`)

	// Output:
	// ERROR: recall statement variable '?next' not found in matches for production 'start' (line 10, col 20)
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
        	recall [foo: ?next *]
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
	// ERROR: could not find chunk named 'foo' (line 8, col 16)
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
		match { goal [foo: * ?other] }
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
			goal [foo: * ?other]
			retrieval [foo: * ?other1]
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
	// ERROR: cannot use ID 'fooID' in print statement (line 9, col 13)
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
	// ERROR: print statement variable '?fooVar' not found in matches for production 'start' (line 9, col 13)
}

func Example_productionPrintStatementWildcard() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	==init==
	==productions==
	start {
		match { retrieval [_status: error] }
		do { print * }
	}`)

	// Output:
	// ERROR: unexpected token "*" (expected "}") (line 9, col 13)
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
	// ERROR: invalid chunk - '_status' expects 1 slot (line 8, col 20)
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
	// ERROR: invalid _status 'something' for 'goal' in production 'start' (should be 'busy', 'empty', 'error', 'full') (line 8, col 25)
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
	// ERROR: invalid _status 'something' for 'retrieval' in production 'start' (should be 'busy', 'empty', 'error', 'full') (line 8, col 30)
}
