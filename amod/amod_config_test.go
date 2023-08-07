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
	// ERROR: unrecognized option "foo" in gactar section (line 5, col 10)
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
	// ERROR: 'log_level' invalid type (found id; expected string) (line 5, col 21)
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
	// ERROR: unrecognized option "foo" in gactar section (line 5, col 10)
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
	// ERROR: 'log_level' invalid type (found field; expected string) (line 5, col 21)
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
	// ERROR: 'trace_activations' invalid type (found number; expected true or false) (line 5, col 29)
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

func Example_modulesExtraBuffers() {
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

func Example_modulesExtraBuffersDuplicate() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules { 
		extra_buffers {
			buffer1 {}
			buffer1 {}
		} 
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: duplicate option "buffer1" (line 8, col 3)
}

func Example_modulesInitBuffer() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		memory {
			max_spread_strength: 0.9
			retrieval { spreading_activation: 0.5 }
		}
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
}

func Example_modulesInitBufferNoMaxSpread() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		imaginal {
			imaginal { spreading_activation: 0.5 }
		}
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: spreading_activation set on buffer "imaginal", but max_spread_strength not set on memory module (line 5, col 1)
}

func Example_modulesInitBufferDuplicate() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		memory {
			max_spread_strength: 0.9
			retrieval { 
				spreading_activation: 0.5
				spreading_activation: 0.5
			}
		}
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: duplicate option "spreading_activation" (line 10, col 4)
}

func Example_modulesUnrecognizedModule() {
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

func Example_modulesUnrecognizedModuleOption() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		goal { foo: 0.2 }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: unrecognized option "foo" in goal config (line 6, col 9)
}

func Example_modulesUnrecognizedModuleDuplicateOption() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		imaginal { 
			imaginal {}
			imaginal {}
		 }
	}
	~~ init ~~
	~~ productions ~~`)

	// Output:
	// ERROR: duplicate option "imaginal" (line 8, col 3)
}

func Example_modulesAll() {
	generateToStdout(`
	~~ model ~~
	name: Test
	~~ config ~~
	modules {
		imaginal {
			delay: 0.2
			imaginal {
				spreading_activation: 0.5
			} 
		}
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
			retrieval {
				spreading_activation: 0.5
			} 
		}
		procedural {
			default_action_time: 0.06
		}
		goal{
			goal { spreading_activation: 0.5 }
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
	// ERROR: imaginal "delay" invalid type (found string; expected number) (line 6, col 20)
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
	// ERROR: imaginal "delay" is out of range (minimum 0) (line 6, col 20)
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
	// ERROR: unrecognized option "foo" in imaginal config (line 6, col 13)
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
	// ERROR: unrecognized option "foo" in memory config (line 6, col 11)
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
	// ERROR: memory "decay" is out of range (0-1) (line 6, col 18)
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
	// ERROR: memory "decay" is out of range (0-1) (line 6, col 18)
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
	// ERROR: unrecognized option "foo" in procedural config (line 6, col 15)
}
