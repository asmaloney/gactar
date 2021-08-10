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

	if set.Field != nil {
		match := production.LookupMatchByBuffer(name)

		if match == nil {
			errs.Addc(&set.Pos, "match buffer '%s' not found in production '%s'", name, production.Name)
		} else {
			if set.Field.ArgNum != nil {
				argNum := int(*set.Field.ArgNum)
				if (argNum == 0) || (argNum > len(match.Pattern.Fields)) {
					errs.Addc(&set.Pos, "field %d does not exist in match buffer '%s' in production '%s'", argNum, name, production.Name)
				}
			} else if set.Field.Name != nil {
				fieldName := *set.Field.Name

				field := match.Pattern.LookupFieldName(fieldName)
				if field == nil {
					errs.Addc(&set.Pos, "field '%s' does not exist in match buffer '%s' in production '%s'", fieldName, name, production.Name)
				}
			} else {
				// should not be possible to get here since the parser will pick this up
				errs.Addc(&set.Pos, "set statement is missing a field number or name in production '%s'", production.Name)
			}
		}
	}

	if set.Pattern == nil && set.Arg == nil {
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

	for _, field := range recall.Pattern.Fields {
		for _, f := range field.Items {

			varName := f.getVar()
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
			errs.Addc(&clear.Pos, "buffer not found in production '%s': '%s'", production.Name, name)
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
		errs.Addc(&write.Pos, "text output not found in production '%s': '%s'", production.Name, name)
	}

	return errs.ErrorOrNil()
}
