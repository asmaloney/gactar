package amod

func Example_gactarUnrecognizedField() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	gactar { foo: bar }
	==init==
	==productions==`)

	// Output:
	// ERROR: unrecognized field in gactar section: 'foo' (line 5, col 10)
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
	// ERROR: cannot use reserved chunk name '_internal' (chunks beginning with '_' are reserved) (line 5, col 11)
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
	// ERROR: duplicate chunk name: 'something' (line 7, col 6)
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
	// ERROR: unrecognized module in config: 'foo' (line 6, col 2)
}

func Example_imaginalFields() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	modules {
		imaginal { delay: 0.2 }
		memory { latency_factor: 0.5 }
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
	// ERROR: imaginal delay 'gack' must be a number (line 6, col 20)
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
	// ERROR: imaginal delay '-0.500000' must be a positive number (line 6, col 20)
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
	// ERROR: unrecognized field 'foo' in imaginal config (line 6, col 13)
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
	// ERROR: unrecognized field 'foo' in memory config (line 6, col 11)
}

func Example_proceduralFieldUnrecognized() {
	generateToStdout(`
	==model==
	name: Test
	==config==
	modules {
		procedural { foo: bar }
	}
	==init==
	==productions==`)

	// Output:
	// ERROR: unrecognized field 'foo' in procedural config (line 6, col 15)
}
