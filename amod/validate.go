package amod

import "gitlab.com/asmaloney/gactar/actr"

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
		} else {
			if set.Slot.ArgNum != nil {
				argNum := int(*set.Slot.ArgNum)
				if (argNum == 0) || (argNum > len(match.Pattern.Slots)) {
					errs.Addc(&set.Pos, "slot %d does not exist in match buffer '%s' in production '%s'", argNum, name, production.Name)
				}
			} else if set.Slot.Name != nil {
				slotName := *set.Slot.Name

				slot := match.Pattern.LookupSlotName(slotName)
				if slot == nil {
					errs.Addc(&set.Pos, "slot '%s' does not exist in match buffer '%s' in production '%s'", slotName, name, production.Name)
				}
			} else {
				// should not be possible to get here since the parser will pick this up
				errs.Addc(&set.Pos, "set statement is missing a slot number or name in production '%s'", production.Name)
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

	name := recall.MemoryName
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
