package amod

func Example_gactarAllOptions() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	gactar {
		log_level: 'detail'
		trace_activations: true
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
}

func Example_gactarUnrecognizedField() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	gactar { foo: bar }
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: unrecognized option in gactar section: 'foo' (line 5, col 10)
}

func Example_gactarUnrecognizedLogLevel() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	gactar { log_level: bar }
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: 'log_level' must be must be one of "min, info, detail" (line 5, col 21)
}

func Example_gactarUnrecognizedNestedValue() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	gactar { foo {} }
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: unrecognized option in gactar section: 'foo' (line 5, col 10)
}

func Example_gactarFieldNotANestedValue() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	gactar { log_level {} }
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: 'log_level' must be must be one of "min, info, detail" (line 5, col 21)
}

func Example_gactarSpaceSeparator() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	gactar { log_level: 'detail' trace_activations: true }
	~~ init ~~
	~~ productions ~~`)

	// Output:
}

func Example_gactarTraceActivations() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	gactar { trace_activations: true }
	~~ init ~~
	~~ productions ~~`)

	// Output:
}

func Example_gactarTraceActivationsNonBool() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	gactar { trace_activations: 6.0 }
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: 'trace_activations' must be 'true' or 'false' (line 5, col 29)
}

func Example_chunkInternalType() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [_internal: foo bar] }
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: cannot use reserved chunk type "_internal" (chunks beginning with '_' are reserved) (line 5, col 11)
}

func Example_chunkReservedType() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks { [requested: foo bar] }
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: cannot use reserved chunk type "requested" (line 5, col 11)
}

func Example_chunkDuplicateType() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	chunks {
    	[something: foo bar]
    	[something: foo bar]
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: duplicate chunk type: 'something' (line 7, col 6)
}

func Example_modules() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		imaginal { delay: 0.2 }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
}

func Example_modulesMultipleBuffers() {
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
	~~ init ~~
	~~ productions ~~`)

	// Output:
}

func Example_modulesUnrecognized() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		foo { delay: 0.2 }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: unrecognized module in config: 'foo' (line 6, col 2)
}

func Example_modulesAll() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		imaginal { delay: 0.2 }
		memory {
			latency_factor: 0.5
			latency_exponent: 0.75
			retrieval_threshold: 0.1
			finst_size: 5
			finst_time: 2.5
			decay: 0.6
			max_spread_strength: 0.9
			instantaneous_noise: 0.5
			mismatch_penalty: 1.0
		}
		procedural {
			default_action_time: 0.06
		}
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	//
}

func Example_imaginalFields() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		imaginal { delay: 0.2 }
		memory { latency_factor: 0.5 }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
}

func Example_imaginalFieldType() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		imaginal { delay: "gack" }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: imaginal 'delay' must be a number (line 6, col 20)
}

func Example_imaginalFieldRange() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		imaginal { delay: -0.5 }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: imaginal 'delay' is out of range (minimum 0) (line 6, col 20)
}

func Example_imaginalFieldUnrecognized() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		imaginal { foo: bar }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: unrecognized option in imaginal config: 'foo' (line 6, col 13)
}

func Example_memoryFieldUnrecognized() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		memory { foo: bar }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: unrecognized option in memory config: 'foo' (line 6, col 11)
}

func Example_memoryDecayOutOfRange() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		memory { decay: -0.5 }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: memory 'decay' is out of range (0-1) (line 6, col 18)
}

func Example_memoryDecayOutOfRange2() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		memory { decay: 1.5 }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: memory 'decay' is out of range (0-1) (line 6, col 18)
}

func Example_proceduralFieldUnrecognized() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		procedural { foo: bar }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: unrecognized option in procedural config: 'foo' (line 6, col 15)
}
