package amod

import (
	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amodlog"
)

// validateChunk checks the chunk name to ensure uniqueness and that it isn't using
// reserved names.
func validateChunk(model *actr.Model, log *amodlog.Log, chunk *chunkDecl) (err error) {
	if actr.IsInternalChunkName(chunk.Name) {
		log.Error(chunk.Pos.Line, "cannot use reserved chunk name '%s' (chunks begining with '_' are reserved)", chunk.Name)
		return CompileError{}
	}

	c := model.LookupChunk(chunk.Name)
	if c != nil {
		log.Error(chunk.Pos.Line, "duplicate chunk name: '%s'", chunk.Name)
		return CompileError{}
	}

	return nil
}

func validateInitialization(model *actr.Model, log *amodlog.Log, init *initialization) (err error) {
	if init.InitPattern != nil {
		err = validatePattern(model, log, init.InitPattern)
	}

	name := init.Name
	buffer := model.LookupBuffer(name)
	if buffer != nil {
		if init.InitPatterns != nil {
			log.Error(init.Pos.Line, "buffer '%s' should only have one pattern in initialization", name)
			return CompileError{}
		} else if init.InitPattern == nil {
			log.Error(init.Pos.Line, "buffer '%s' requires a pattern in initialization", name)
			return CompileError{}
		}

		return
	}

	if name != "memory" {
		log.Error(init.Pos.Line, "buffer or memory '%s' not found in initialization ", name)
		err = CompileError{}
	}

	for _, init := range init.InitPatterns {
		pattern_err := validatePattern(model, log, init)
		if pattern_err != nil {
			err = CompileError{}
			continue
		}
	}

	return
}

// validatePattern ensures that the pattern's chunk exists and that its number of slots match.
func validatePattern(model *actr.Model, log *amodlog.Log, pattern *pattern) (err error) {
	chunkName := pattern.ChunkName
	chunk := model.LookupChunk(chunkName)
	if chunk == nil {
		log.Error(pattern.Pos.Line, "could not find chunk named '%s'", chunkName)
		return CompileError{}
	}

	if len(pattern.Slots) != chunk.NumSlots {
		log.Error(pattern.Pos.Line, "invalid chunk - '%s' expects %d slots", chunkName, chunk.NumSlots)
		return CompileError{}
	}

	return
}

// validateMatch verifies several aspects of a match item.
func validateMatch(match *match, model *actr.Model, log *amodlog.Log, production *actr.Production) (err error) {
	if match == nil {
		return
	}

	for _, item := range match.Items {
		name := item.Name

		buffer := model.LookupBuffer(name)

		if buffer == nil {
			log.Error(item.Pos.Line, "buffer '%s' not found in production '%s'", name, production.Name)
			err = CompileError{}
			continue
		}

		pattern := item.Pattern
		if pattern == nil {
			log.Error(item.Pos.Line, "invalid pattern for '%s' in production '%s'", name, production.Name)
			err = CompileError{}
			continue
		}

		// check _status chunks to ensure they have one of the allowed tests
		if pattern.ChunkName == "_status" {
			if len(pattern.Slots) != 1 {
				log.Error(item.Pos.Line, "_status should only have one slot for '%s' in production '%s' (should be %s)", name, production.Name, actr.ValidBufferStatesStr())
				err = CompileError{}
			}

			slot := *pattern.Slots[0].Items[0].ID

			if !actr.IsValidBufferState(slot) {
				log.Error(item.Pos.Line, "invalid _status '%s' for '%s' in production '%s' (should be %v)", slot, name, production.Name, actr.ValidBufferStatesStr())
				err = CompileError{}
			}
		}
	}

	return
}

// validateSetStatement checks a "set" statement to verify the buffer name & field indexing is correct.
// The production's matches have been constructed, so that's what we check against.
func validateSetStatement(set *setStatement, model *actr.Model, log *amodlog.Log, production *actr.Production) (err error) {
	bufferName := set.BufferName
	buffer := model.LookupBuffer(bufferName)
	if buffer == nil {
		log.Error(set.Pos.Line, "buffer '%s' not found in production '%s'", bufferName, production.Name)
		err = CompileError{}
	}

	if set.Slot != nil {
		// we should have the form "set <slot name> of <buffer> to <value>"
		slotName := *set.Slot
		if set.Pattern != nil {
			log.Error(set.Pos.Line, "cannot set a slot ('%s') to a pattern in match buffer '%s' in production '%s'", slotName, bufferName, production.Name)
			err = CompileError{}
		} else {
			match := production.LookupMatchByBuffer(bufferName)

			if match == nil {
				log.Error(set.Pos.Line, "match buffer '%s' not found in production '%s'", bufferName, production.Name)
				err = CompileError{}
			} else {
				chunk := match.Pattern.Chunk
				if chunk == nil {
					log.Error(set.Pos.Line, "chunk does not exist in match buffer '%s' in production '%s'", bufferName, production.Name)
					err = CompileError{}
				} else {
					if !chunk.HasSlot(slotName) {
						log.Error(set.Pos.Line, "slot '%s' does not exist in chunk '%s' for match buffer '%s' in production '%s'", slotName, chunk.Name, bufferName, production.Name)
						err = CompileError{}
					}

					if set.Value.Var != nil {
						// Check set.Value.Var to ensure it exists
						match := production.LookupMatchByVariable(*set.Value.Var)
						if match == nil {
							log.Error(set.Value.Pos.Line, "set statement variable '%s' not found in matches for production '%s'", *set.Value.Var, production.Name)
							err = CompileError{}
						}
					}
				}
			}
		}
	} else {
		// we should have the form "set <buffer> to <pattern>"
		if set.Value != nil {
			log.Error(set.Pos.Line, "buffer '%s' must be set to a pattern in production '%s'", bufferName, production.Name)
			err = CompileError{}
		}
	}

	if set.Pattern == nil && set.Value == nil {
		// should not be possible to get here since the parser should pick this up
		log.Error(set.Pos.Line, "set statement is missing value (set to what?) in production '%s'", production.Name)
		err = CompileError{}
	}

	return
}

// validateRecallStatement checks a "recall" statement to verify the memory name.
func validateRecallStatement(recall *recallStatement, model *actr.Model, log *amodlog.Log, production *actr.Production) (err error) {
	for _, slot := range recall.Pattern.Slots {
		for _, item := range slot.Items {
			varName := item.Var
			if (varName == nil) || (*varName == "?") {
				continue
			}

			match := production.LookupMatchByVariable(*varName)
			if match == nil {
				log.Error(recall.Pos.Line, "recall statement variable '%s' not found in matches for production '%s'", *varName, production.Name)
				err = CompileError{}
			}
		}
	}

	return
}

// validateClearStatement checks a "clear" statement to verify the buffer names.
func validateClearStatement(clear *clearStatement, model *actr.Model, log *amodlog.Log, production *actr.Production) (err error) {
	bufferNames := clear.BufferNames

	for _, name := range bufferNames {
		buffer := model.LookupBuffer(name)
		if buffer == nil {
			log.Error(clear.Pos.Line, "buffer '%s' not found in production '%s'", name, production.Name)

			err = CompileError{}
			continue
		}
	}

	return
}

// validatePrintStatement is a placeholder for checking a "print" statement. Currently there are no checks.
func validatePrintStatement(print *printStatement, model *actr.Model, log *amodlog.Log, production *actr.Production) (err error) {
	if print.Args != nil {
		for _, v := range print.Args {
			if v.ID != nil {
				log.Error(print.Pos.Line, "cannot use ID '%s' in print statement", *v.ID)
			} else if v.Var != nil {
				match := production.LookupMatchByVariable(*v.Var)
				if match == nil {
					log.Error(print.Pos.Line, "print statement variable '%s' not found in matches for production '%s'", *v.Var, production.Name)
					err = CompileError{}
				}
			}
		}
	}

	return
}
