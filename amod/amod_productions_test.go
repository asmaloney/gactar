package amod

func Example_production() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?blat] }
		do {
			print ?blat
			stop
		}
	}`)

	// Output:
}

func Example_productionNoDo() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?blat] }
	}`)

	// Output:
	// ERROR: unexpected token "}" (expected Do "}") (line 10, col 1)
}

func Example_productionWhenClause() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match {
			goal [foo: ?blat] when ( ?blat == 'foo' )
		}
		do {
			print ?blat
			stop
		}
	}`)

	// Output:
}

func Example_productionWhenClauseCompareNil() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match {
			goal [foo: ?blat] when ( ?blat == nil )
		}
		do {
			print ?blat
			stop
		}
	}`)

	// Output:
}

func Example_productionWhenClauseCompareID() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match {
			goal [foo: ?blat] when ( ?blat == bar )
		}
		do {
			print ?blat
			stop
		}
	}`)

	// Output:
	// ERROR: unexpected token "bar" (expected WhenArg ")") (line 10, col 37)
}

func Example_productionWhenClauseNegatedAndConstrained() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match {
			goal [foo: !?blat]
				when ( ?blat == 'foo' )
		}
		do {
			print ?blat
			stop
		}
	}`)

	// Output:
	// ERROR: cannot further constrain a negated variable '?blat' (line 11, col 11)
}

func Example_productionWhenClauseComparisonToSelf() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match {
			goal [foo: ?blat] when ( ?blat == ?blat )
		}
		do {
			print ?blat
			stop
		}
	}`)

	// Output:
	// ERROR: cannot compare a variable to itself '?blat' (line 10, col 37)
}

func Example_productionWhenClauseInvalidVarLHS() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match {
			goal [foo: ?blat] when ( ?ding == 42 )
		}
		do {
			print ?blat
			stop
		}
	}`)

	// Output:
	// ERROR: unknown variable ?ding in where clause (line 10, col 28)
}

func Example_productionWhenClauseInvalidVarRHS() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
	start {
		match {
			goal [foo: ?blat] when ( ?blat != ?ding )
		}
		do {
			print ?blat
			stop
		}
	}`)

	// Output:
	// ERROR: unknown variable ?ding in where clause (line 10, col 37)
}

func Example_productionWildcard() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
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
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?blat !*] }
		do { print ?blat }
	}`)

	// Output:
	// ERROR: unexpected token "!" (expected "]") (line 9, col 27)
}

func Example_productionUnusedVar1() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing] }
	~~ init ~~
	~~ productions ~~
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
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [foo: thing1 thing2] }
	~~ init ~~
	~~ productions ~~
	start {
		match { goal [foo: ?blat ?blat] }
		do { set goal to [foo: ding] }
	}`)

	// Output:
}

func Example_productionInvalidMemory() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { another_goal [add: ?thing1 ?thing2] }
		do { print 'foo' }
	}`)

	// Output:
	// ERROR: buffer 'another_goal' not found in production 'start' (line 8, col 10)
}

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

func Example_productionClearStatementInvalidBuffer() {
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

func Example_productionSetStatementNonBuffer() {
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

func Example_productionSetStatementNonBuffer2() {
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

func Example_productionSetStatementInvalidSlot() {
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

func Example_productionSetStatementNonVar1() {
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

func Example_productionSetStatementNonVar2() {
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

func Example_productionSetStatementAssignNonPattern() {
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

func Example_productionSetStatementAssignPattern() {
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

func Example_productionRecallStatementMultiple() {
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

func Example_productionRecallStatementInvalidPattern() {
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

func Example_productionRecallStatementVarNotFound() {
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

func Example_productionRecallStatementWithInvalidParam() {
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

func Example_productionRecallStatementWithString() {
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

func Example_productionRecallStatementWithVar() {
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

func Example_productionRecallStatementWithIncorrectNumArgs() {
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

func Example_productionChunkNotFound() {
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

func Example_productionPrintStatementInvalidBuffer() {
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

func Example_productionPrintStatementInvalidBufferSlot() {
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

func Example_productionPrintStatementInvalidNil() {
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

func Example_productionPrintStatementInvalidVar() {
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

func Example_productionPrintStatementWildcard() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { buffer_state goal empty }
		do { print * }
	}`)

	// Output:
	// ERROR: unexpected token "*" (expected "}") (line 9, col 13)
}

func Example_productionMatchBufferState() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { buffer_state retrieval empty }
		do { print 42 }
	}`)

	// Output:
}

func Example_productionMatchBufferStateInvalidBuffer() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { buffer_state foo empty }
		do { print 42 }
	}`)

	// Output:
	// ERROR: buffer 'foo' not found in production 'start' (line 8, col 10)
}

func Example_productionMatchBufferStateInvalidStatus() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { buffer_state retrieval foo }
		do { print 42 }
	}`)

	// Output:
	// ERROR: invalid state check 'foo' for buffer 'retrieval' in production 'start' (should be one of: empty, full) (line 8, col 10)
}

func Example_productionMatchBufferStateInvalidString() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { buffer_state goal 'empty' }
		do { print 42 }
	}`)

	// Output:
	// ERROR: unexpected token "empty" (expected <ident>) (line 8, col 28)
}

func Example_productionMatchBufferStateInvalidNumber() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { buffer_state goal 42 }
		do { print 42 }
	}`)

	// Output:
	// ERROR: unexpected token "42" (expected <ident>) (line 8, col 28)
}

func Example_productionMatchBufferStateDuplicate() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match {
			buffer_state retrieval full
			buffer_state retrieval empty
		}
		do { print 42 }
	}`)

	// Output:
	// ERROR: duplicate buffer state check for 'retrieval' in production 'start' (line 10, col 3)
}

func Example_productionMatchModuleState() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { module_state memory error }
		do { print 42 }
	}`)

	// Output:
}

func Example_productionMatchModuleStateInvalidModule() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { module_state foo error }
		do { print 42 }
	}`)

	// Output:
	// ERROR: module 'foo' not found in production 'start' (line 8, col 10)
}

func Example_productionMatchModuleStateInvalidState1() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { module_state memory 'foo' }
		do { print 42 }
	}`)

	// Output:
	// ERROR: unexpected token "foo" (expected <ident>) (line 8, col 30)
}

func Example_productionMatchModuleStateInvalidState2() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { module_state memory bar }
		do { print 42 }
	}`)

	// Output:
	// ERROR: invalid module state check 'bar' for module 'memory' in production 'start' (should be one of: busy, error, free) (line 8, col 10)
}

func Example_productionMatchModuleStateDuplicate() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { 
			module_state memory error
			module_state memory busy
		}
		do { print 42 }
	}`)

	// Output:
	// ERROR: duplicate module state check for 'memory' in production 'start' (line 10, col 3)
}
