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
		[author: 'me' 'software']
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
	goal [author: 'Fred' 'Book' '1972']
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
	memory [author: 'Jane' 'Book' '1982']
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
			[remember: 'me']
			[author: 'me' 'software']
		}
	}
	~~ productions ~~`)

	// Output:
}

func Example_initializer5() {
	// memory with named chunks
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [author: person object year] }
	~~ init ~~
	memory {
		bar [author: 'Fred' 'Book' '1972']
		foo [author: 'Jane' 'Book' '1982']
		[author: 'Xe' 'Software' '2001']
	}
	~~ productions ~~`)

	// Output:
}

func Example_initializer6() {
	// memory with one init and named chunk
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [author: person object year] }
	~~ init ~~
	memory foo [author: 'Jane' 'Book' '1982']
	~~ productions ~~`)

	// Output:
}

func Example_initializerErrorInvalidSlots() {
	// Check invalid number of slots in init
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [author: person object year] }
	~~ init ~~
	memory { [author: 'me' 'software'] }
	~~ productions ~~`)

	// Output:
	// ERROR: invalid chunk - 'author' expects 3 slots (line 7, col 10)
}

func Example_initializerErrorInvalidChunk1() {
	// Check memory with invalid chunk
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	memory { [author: 'me' 'software'] }
	~~ productions ~~`)

	// Output:
	// ERROR: could not find chunk named 'author' (line 6, col 11)
}

func Example_initializerErrorInvalidChunk2() {
	// Check buffer with invalid chunk
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	goal [author: 'Fred' 'Book' '1972']
	~~ productions ~~`)

	// Output:
	// ERROR: could not find chunk named 'author' (line 6, col 7)
}

func Example_initializerErrorUnknownBuffer() {
	// Check unknown buffer
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [author: person object year] }
	~~ init ~~
	something [author: 'Fred' 'Book' '1972']
	~~ productions ~~`)

	// Output:
	// ERROR: module 'something' not found in initialization (line 7, col 1)
}

func Example_initializerErrorMultipleInits() {
	// Check buffer with multiple inits
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [author: person object year] }
	~~ init ~~
	goal { [author: 'Fred' 'Book' '1972'] [author: 'Jane' 'Book' '1982'] }
	~~ productions ~~`)

	// Output:
	// ERROR: module "goal" should only have one pattern in initialization of buffer "goal" (line 7, col 8)
}

func Example_initializerMultipleBuffers() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules { 
		extra_buffers {
			buffer1 {}
			buffer2 {}
		} 
	}
	chunks { [author: person object year] }
	~~ init ~~
	extra_buffers {
		buffer1 { [author: 'Fred' 'Book' '1972'] }
		buffer2 { [author: 'Jane' 'Book' '1984'] }
	}
	~~ productions ~~`)

	// Output:
}

func Example_initializerErrorNoBuffers() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules { extra_buffers{} }
	chunks { [author: person object year] }
	~~ init ~~
	extra_buffers {
		[author: 'Jane' 'Book' '1984']
		[author: 'Xe' 'Book' '1999']
	}
	~~ productions ~~`)

	// Output:
	// ERROR: module 'extra_buffers' does not have any buffers (line 8, col 1)
}

func Example_initializerErrorDuplicateNames() {
	// memory with named chunks
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [author: person object year] }
	~~ init ~~
	memory {
		foo [author: 'Fred' 'Book' '1972']
		foo [author: 'Jane' 'Book' '1982']
		[author: 'Xe' 'Software' '2001']
	}
	~~ productions ~~`)

	// Output:
	// ERROR: duplicate chunk name "foo" found in initialization (line 9, col 2)
}

func Example_initializerPartialSimilarities() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		memory {
			mismatch_penalty: 1.0
		}
	}
	chunks {
    	[group: id parent position]
	}
	~~ init ~~
	memory {
		[group: group1 list first]
		[group: group2 list second]
		[group: group3 list third]
	}

	similar {
		( first second -0.5 )
		( second third -0.5 )
	}

	~~ productions ~~`)

	// Output:
}
