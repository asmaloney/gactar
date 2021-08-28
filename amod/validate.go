package amod

import (
	"gitlab.com/asmaloney/gactar/actr"
)

// validateChunk checks the chunk name to ensure uniqueness and that it isn't using
// reserved names.
func validateChunk(model *actr.Model, chunk *chunkDecl) (err error) {
	errs := errorListWithContext{}

	if actr.IsInternalChunkName(chunk.Name) {
		errs.Addc(&chunk.Pos, "cannot use reserved chunk name '%s' (chunks begining with '_' are reserved)", chunk.Name)
		return errs
	}

	if actr.ReservedChunkNameExists(chunk.Name) {
		errs.Addc(&chunk.Pos, "cannot use reserved chunk name '%s'", chunk.Name)
		return errs
	}

	c := model.LookupChunk(chunk.Name)
	if c != nil {
		errs.Addc(&chunk.Pos, "duplicate chunk name: '%s'", chunk.Name)
	}

	return errs.ErrorOrNil()
}

// validatePattern ensures that the pattern's chunk exists and that its number of slots match.
func validatePattern(model *actr.Model, pattern *pattern) (err error) {
	errs := errorListWithContext{}

	chunkName := pattern.ChunkName
	chunk := model.LookupChunk(chunkName)
	if chunk == nil {
		errs.Addc(&pattern.Pos, "could not find chunk named '%s'", chunkName)
		return errs
	}

	if len(pattern.Slots) != chunk.NumSlots {
		errs.Addc(&pattern.Pos, "invalid chunk - '%s' expects %d slots", chunkName, chunk.NumSlots)
	}

	return errs.ErrorOrNil()
}

// validateMatch verifies several aspects of a match item.
func validateMatch(match *match, model *actr.Model, production *actr.Production) (err error) {
	errs := errorListWithContext{}

	if match == nil {
		return
	}

	for _, item := range match.Items {
		name := item.Name

		buffer := model.LookupBuffer(name)
		memory := model.LookupMemory(name)

		if (buffer == nil) && (memory == nil) {
			errs.Addc(&item.Pos, "buffer or memory '%s' not found in production '%s'", name, production.Name)
			continue
		}

		pattern := item.Pattern
		if pattern == nil {
			errs.Addc(&item.Pos, "invalid pattern for '%s' in production '%s'", name, production.Name)
			continue
		}

		// check _status chunks to ensure they have one of the allowed tests
		if buffer != nil {
			if pattern.ChunkName == "_status" {
				if len(pattern.Slots) != 1 {
					errs.Addc(&item.Pos, "_status should only have one slot for '%s' in production '%s' (should be 'full' or 'empty')", name, production.Name)
				}

				slot := *pattern.Slots[0].Items[0].ID

				if slot != "full" && slot != "empty" {
					errs.Addc(&item.Pos, "invalid _status '%s' for '%s' in production '%s' (should be 'full' or 'empty')", slot, name, production.Name)
				}
			}
		} else if memory != nil {
			if pattern.ChunkName == "_status" {
				if len(pattern.Slots) != 1 {
					errs.Addc(&item.Pos, "_status should only have one slot for '%s' in production '%s' (should be 'busy', 'free', or 'error')", name, production.Name)
				}

				slot := *pattern.Slots[0].Items[0].ID

				if slot != "busy" && slot != "free" && slot != "error" {
					errs.Addc(&item.Pos, "invalid _status '%s' for '%s' in production '%s' (should be 'busy', 'free', or 'error')", slot, name, production.Name)
				}
			}
		}
	}

	return errs.ErrorOrNil()
}

// validateSetStatement checks a "set" statement to verify the buffer name & field indexing is correct.
// The production's matches have been constructed, so that's what we check against.
func validateSetStatement(set *setStatement, model *actr.Model, production *actr.Production) (err error) {
	errs := errorListWithContext{}

	bufferName := set.BufferName
	buffer := model.LookupBuffer(bufferName)
	if buffer == nil {
		errs.Addc(&set.Pos, "buffer '%s' not found in production '%s'", bufferName, production.Name)
	}

	if set.Slot != nil {
		// we should have the form "set <slot name> of <buffer> to <value>"
		slotName := *set.Slot
		if set.Pattern != nil {
			errs.Addc(&set.Pos, "cannot set a slot ('%s') to a pattern in match buffer '%s' in production '%s'", slotName, bufferName, production.Name)
		} else {
			match := production.LookupMatchByBuffer(bufferName)

			if match == nil {
				errs.Addc(&set.Pos, "match buffer '%s' not found in production '%s'", bufferName, production.Name)
			} else {
				chunk := match.Pattern.Chunk
				if chunk == nil {
					errs.Addc(&set.Pos, "chunk does not exist in match buffer '%s' in production '%s'", bufferName, production.Name)
				} else {
					if !chunk.SlotExists(slotName) {
						errs.Addc(&set.Pos, "slot '%s' does not exist in chunk '%s' for match buffer '%s' in production '%s'", slotName, chunk.Name, bufferName, production.Name)
					}

					if set.Value.Var != nil {
						// Check set.Value.Var to ensure it exists
						match := production.LookupMatchByVariable(*set.Value.Var)
						if match == nil {
							errs.Addc(&set.Value.Pos, "set statement variable '%s' not found in matches for production '%s'", *set.Value.Var, production.Name)
						}
					}
				}
			}
		}
	} else {
		// we should have the form "set <buffer> to <pattern>"
		if set.Value != nil {
			errs.Addc(&set.Pos, "buffer '%s' must be set to a pattern in production '%s'", bufferName, production.Name)
		}
	}

	if set.Pattern == nil && set.Value == nil {
		// should not be possible to get here since the parser should pick this up
		errs.Addc(&set.Pos, "set statement is missing value (set to what?) in production '%s'", production.Name)
	}

	return errs.ErrorOrNil()
}

// validateRecallStatement checks a "recall" statement to verify the memory name.
func validateRecallStatement(recall *recallStatement, model *actr.Model, production *actr.Production) (err error) {
	errs := errorListWithContext{}

	name := "memory"
	memory := model.LookupMemory(name)
	if memory == nil {
		errs.Addc(&recall.Pos, "recall statement memory '%s' not found in production '%s'", name, production.Name)
	}

	for _, slot := range recall.Pattern.Slots {
		for _, item := range slot.Items {
			varName := item.Var
			if (varName == nil) || (*varName == "?") {
				continue
			}

			match := production.LookupMatchByVariable(*varName)
			if match == nil {
				errs.Addc(&recall.Pos, "recall statement variable '%s' not found in matches for production '%s'", *varName, production.Name)
			}
		}
	}

	return errs.ErrorOrNil()
}

// validateClearStatement checks a "clear" statement to verify the buffer names.
func validateClearStatement(clear *clearStatement, model *actr.Model, production *actr.Production) (err error) {
	errs := errorListWithContext{}

	bufferNames := clear.BufferNames

	for _, name := range bufferNames {
		buffer := model.LookupBuffer(name)
		if buffer == nil {
			errs.Addc(&clear.Pos, "buffer '%s' not found in production '%s'", name, production.Name)
			continue
		}
	}

	return errs.ErrorOrNil()
}

// validatePrintStatement is a placeholder for checking a "print" statement. Currently there are no checks.
func validatePrintStatement(print *printStatement, model *actr.Model, production *actr.Production) (err error) {
	errs := errorListWithContext{}

	if print.Args != nil {
		for _, v := range print.Args {
			if v.ID != nil {
				errs.Addc(&print.Pos, "cannot use ID '%s' in print statement", *v.ID)
			} else if v.Var != nil {
				match := production.LookupMatchByVariable(*v.Var)
				if match == nil {
					errs.Addc(&print.Pos, "print statement variable '%s' not found in matches for production '%s'", *v.Var, production.Name)
				}
			}
		}
	}

	return errs.ErrorOrNil()
}
