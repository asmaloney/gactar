package amod

import (
	"gitlab.com/asmaloney/gactar/actr"
)

// validateInitializers looks at the initializers and ensures that the number of slots is valid relative
// to the number of slot names. Note this has to happen after we've determined the slot names.
func validateInitializers(model *actr.Model, init *initSection) (err error) {
	errs := errorListWithContext{}

	if init == nil {
		return
	}

	for _, i := range init.Initializers {
		memory := model.LookupMemory("memory")
		if memory == nil {
			errs.Addc(&i.Pos, "memory not found")
			continue
		}

		for line, str := range i.Items.Strings {
			// we need to guess the line number since we just have an array of strings here
			pos := i.Pos
			pos.Line += line + 1

			chunkName, slots := actr.SplitStringForChunk(str)

			chunk := model.LookupChunk(chunkName)
			if chunk == nil {
				errs.Addc(&pos, "could not find chunk named '%s' in initialization '%s' for memory (line number approximate)", chunkName, str)
				continue
			}

			if len(slots) != chunk.NumSlots {
				errs.Addc(&pos, "invalid initialization '%s' for memory - expected %d slots (line number approximate)", str, chunk.NumSlots)
				continue
			}
		}
	}

	return errs.ErrorOrNil()
}

// validateSetStatement checks a "set" statement to verify the buffer name & field indexing is correct.
// The production's matches have been constructed, so that's what we check against.
func validateSetStatement(set *setStatement, model *actr.Model, production *actr.Production) (err error) {
	errs := errorListWithContext{}

	name := set.BufferName
	buffer := model.LookupBuffer(name)
	if buffer == nil {
		errs.Addc(&set.Pos, "buffer '%s' not found in production '%s'", name, production.Name)
	}

	if set.Slot != nil {
		match := production.LookupMatchByBuffer(name)

		if match == nil {
			errs.Addc(&set.Pos, "match buffer '%s' not found in production '%s'", name, production.Name)
		} else if set.Slot == nil {
			// should not be possible to get here since the parser will pick this up
			errs.Addc(&set.Pos, "set statement is missing a slot number or name in production '%s'", production.Name)
		} else {
			slotName := *set.Slot
			chunk := match.Pattern.Chunk
			if chunk == nil {
				errs.Addc(&set.Pos, "chunk does not exist in match buffer '%s' in production '%s'", name, production.Name)
			} else {
				if !chunk.SlotExists(slotName) {
					errs.Addc(&set.Pos, "slot '%s' does not exist in chunk '%s' for match buffer '%s' in production '%s'", slotName, chunk.Name, name, production.Name)
				}
			}
		}
	}

	if set.Pattern == nil && set.ID == nil && set.Number == nil && set.String == nil {
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

			varName := item.getVar()
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

	return errs.ErrorOrNil()
}

// validateWriteStatement checks a "write" statement to verify the text output name.
func validateWriteStatement(write *writeStatement, model *actr.Model, production *actr.Production) (err error) {
	errs := errorListWithContext{}

	name := write.TextOutputName
	textOutput := model.LookupTextOutput(name)
	if textOutput == nil {
		errs.Addc(&write.Pos, "text output '%s' not found in production '%s'", name, production.Name)
	}

	return errs.ErrorOrNil()
}
