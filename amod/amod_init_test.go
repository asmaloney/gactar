package amod

func Example_initializer1() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks {
		[remember: person]
		[author: person object]
	}
	~~ init ~~
	memory {
		[remember: me]
		[author: me software]
	}
	~~ productions ~~`)

	// Output:
}

func Example_initializer2() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [author: person object year] }
	~~ init ~~
	goal [author: Fred Book 1972]
	~~ productions ~~`)

	// Output:
}

func Example_initializer3() {
	// memory with one init is allowed
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [author: person object year] }
	~~ init ~~
	memory [author: Jane Book 1982]
	~~ productions ~~`)

	// Output:
}

func Example_initializer4() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks {
		[remember: person]
		[author: person object]
	}
	~~ init ~~
	memory {
		retrieval {
			[remember: me]
			[author: me software]
		}
	}
	~~ productions ~~`)

	// Output:
}

func Example_initializerInvalidSlots() {
	// Check invalid number of slots in init
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [author: person object year] }
	~~ init ~~
	memory { [author: me software] }
	~~ productions ~~`)

	// Output:
	// ERROR: invalid chunk - 'author' expects 3 slots (line 7, col 10)
}

func Example_initializerInvalidChunk1() {
	// Check memory with invalid chunk
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	memory { [author: me software] }
	~~ productions ~~`)

	// Output:
	// ERROR: could not find chunk named 'author' (line 6, col 11)
}

func Example_initializerInvalidChunk2() {
	// Check buffer with invalid chunk
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	goal [author: Fred Book 1972]
	~~ productions ~~`)

	// Output:
	// ERROR: could not find chunk named 'author' (line 6, col 7)
}

func Example_initializerUnknownBuffer() {
	// Check unknown buffer
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [author: person object year] }
	~~ init ~~
	something [author: Fred Book 1972]
	~~ productions ~~`)

	// Output:
	// ERROR: module 'something' not found in initialization (line 7, col 1)
}

func Example_initializerMultipleInits() {
	// Check buffer with multiple inits
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [author: person object year] }
	~~ init ~~
	goal { [author: Fred Book 1972] [author: Jane Book 1982] }
	~~ productions ~~`)

	// Output:
	// ERROR: module "goal" should only have one pattern in initialization of buffer "goal" (line 7, col 8)
}
