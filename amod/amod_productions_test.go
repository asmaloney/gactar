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

func Example_productionMatchBufferAny() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~
	start {
		match { retrieval [any] }
		do { print 42 }
	}`)

	// Output:
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
