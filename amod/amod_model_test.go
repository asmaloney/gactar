package amod

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
	// ERROR: could not find chunk named 'foo' (line 4, col 13)
}
