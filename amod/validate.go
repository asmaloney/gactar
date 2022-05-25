package amod

import (
	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/issues"
)

// validateChunk checks the chunk name to ensure uniqueness and that it isn't using
// reserved names.
func validateChunk(model *actr.Model, log *issues.Log, chunk *chunkDecl) (err error) {
	if actr.IsInternalChunkName(chunk.Name) {
		log.Error(chunk.Pos.Line, chunk.Pos.Column, "cannot use reserved chunk name '%s' (chunks begining with '_' are reserved)", chunk.Name)
		return CompileError{}
	}

	c := model.LookupChunk(chunk.Name)
	if c != nil {
		log.Error(chunk.Pos.Line, chunk.Pos.Column, "duplicate chunk name: '%s'", chunk.Name)
		return CompileError{}
	}

	return nil
}

func validateInitialization(model *actr.Model, log *issues.Log, init *initialization) (err error) {
	name := init.Name
	module := model.LookupModule(name)

	if module == nil {
		log.Error(init.Pos.Line, init.Pos.Column, "module '%s' not found in initialization ", name)
		return CompileError{}
	}

	if !module.AllowsMultipleInit() && len(init.InitPatterns) > 1 {
		log.Error(init.Pos.Line, init.Pos.Column, "module '%s' should only have one pattern in initialization", name)
		return CompileError{}
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
func validatePattern(model *actr.Model, log *issues.Log, pattern *pattern) (err error) {
	chunkName := pattern.ChunkName
	chunk := model.LookupChunk(chunkName)
	if chunk == nil {
		log.Error(pattern.Pos.Line, pattern.Pos.Column, "could not find chunk named '%s'", chunkName)
		return CompileError{}
	}

	if len(pattern.Slots) != chunk.NumSlots {
		log.Error(pattern.Pos.Line, pattern.Pos.Column, "invalid chunk - '%s' expects %d slots", chunkName, chunk.NumSlots)
		return CompileError{}
	}

	return
}

// validateMatch verifies several aspects of a match item.
func validateMatch(match *match, model *actr.Model, log *issues.Log, production *actr.Production) (err error) {
	if match == nil {
		return
	}

	for _, item := range match.Items {
		name := item.Name

		buffer := model.LookupBuffer(name)

		if buffer == nil {
			log.Error(item.Pos.Line, item.Pos.Column, "buffer '%s' not found in production '%s'", name, production.Name)
			err = CompileError{}
			continue
		}

		pattern := item.Pattern
		if pattern == nil {
			log.Error(item.Pos.Line, item.Pos.Column, "invalid pattern for '%s' in production '%s'", name, production.Name)
			err = CompileError{}
			continue
		}

		// check _status chunks to ensure they have one of the allowed tests
		if pattern.ChunkName == "_status" {
			if len(pattern.Slots) != 1 {
				log.Error(item.Pos.Line, item.Pos.Column, "_status should only have one slot for '%s' in production '%s' (should be %s)", name, production.Name, actr.ValidBufferStatesStr())
				err = CompileError{}
			}

			slot := *pattern.Slots[0].Items[0].ID

			if !actr.IsValidBufferState(slot) {
				log.Error(item.Pos.Line, item.Pos.Column, "invalid _status '%s' for '%s' in production '%s' (should be %v)", slot, name, production.Name, actr.ValidBufferStatesStr())
				err = CompileError{}
			}
		}
	}

	return
}

// validateDo checks for multiple recall statements.
func validateDo(log *issues.Log, production *production) {
	type ref struct {
		lastLine int // keep track of the last case of a recall statement for our output
		count    int // ref count
	}

	recallRef := ref{
		lastLine: 0,
		count:    0,
	}

	for _, statement := range *production.Do.Statements {
		if statement.Recall != nil {
			recallRef.lastLine = statement.Pos.Line
			recallRef.count++
		}
	}

	if recallRef.count > 1 {
		log.Error(recallRef.lastLine, 0, "only one recall statement per production is allowed in production '%s'", production.Name)
	}
}

// validateSetStatement checks a "set" statement to verify the buffer name & field indexing is correct.
// The production's matches have been constructed, so that's what we check against.
func validateSetStatement(set *setStatement, model *actr.Model, log *issues.Log, production *actr.Production) (err error) {
	bufferName := set.BufferName
	buffer := model.LookupBuffer(bufferName)
	if buffer == nil {
		log.Error(set.Pos.Line, set.Pos.Column, "buffer '%s' not found in production '%s'", bufferName, production.Name)
		err = CompileError{}
	}

	if set.Slot != nil {
		// we have the form "set <buffer>.<slot name> to <value>"
		slotName := *set.Slot
		if set.Pattern != nil {
			log.Error(set.Pos.Line, set.Pos.Column, "cannot set a slot ('%s') to a pattern in match buffer '%s' in production '%s'", slotName, bufferName, production.Name)
			err = CompileError{}
		} else {
			match := production.LookupMatchByBuffer(bufferName)

			if match == nil {
				log.Error(set.Pos.Line, set.Pos.Column, "match buffer '%s' not found in production '%s'", bufferName, production.Name)
				err = CompileError{}
			} else {
				chunk := match.Pattern.Chunk
				if chunk == nil {
					log.Error(set.Pos.Line, set.Pos.Column, "chunk does not exist in match buffer '%s' in production '%s'", bufferName, production.Name)
					err = CompileError{}
				} else {
					if !chunk.HasSlot(slotName) {
						log.Error(set.Pos.Line, set.Pos.Column, "slot '%s' does not exist in chunk '%s' for match buffer '%s' in production '%s'", slotName, chunk.Name, bufferName, production.Name)
						err = CompileError{}
					}

					if set.Value.Var != nil {
						// Check set.Value.Var to ensure it exists
						varItem := *set.Value.Var
						match := production.LookupMatchByVariable(varItem)
						if match == nil {
							if varItem == "?" {
								log.Error(set.Value.Pos.Line, set.Value.Pos.Column, "cannot set '%s.%s' to anonymous var ('?') in production '%s'", bufferName, slotName, production.Name)
							} else {
								log.Error(set.Value.Pos.Line, set.Value.Pos.Column, "set statement variable '%s' not found in matches for production '%s'", varItem, production.Name)
							}
							err = CompileError{}
						}
					}
				}
			}
		}
	} else {
		// we have the form "set <buffer> to <pattern>"
		if set.Value != nil {
			log.Error(set.Pos.Line, set.Pos.Column, "buffer '%s' must be set to a pattern in production '%s'", bufferName, production.Name)
			err = CompileError{}
		} else {
			chunkName := set.Pattern.ChunkName
			chunk := model.LookupChunk(chunkName)

			for slotIndex, slot := range set.Pattern.Slots {
				if len(slot.Items) > 1 {
					log.Error(set.Pattern.Pos.Line, set.Pattern.Pos.Column, "cannot set '%s.%v' to compound var in production '%s'", bufferName, chunk.SlotName(slotIndex), production.Name)
					err = CompileError{}

					continue
				}

				// we only have one item
				item := slot.Items[0]
				if item.Var == nil {
					continue
				}

				varItem := *item.Var
				match := production.LookupMatchByVariable(varItem)
				if match == nil {
					if varItem == "?" {
						log.Error(set.Pattern.Pos.Line, set.Pattern.Pos.Column, "cannot set '%s.%v' to anonymous var ('?') in production '%s'", bufferName, chunk.SlotName(slotIndex), production.Name)
					} else {
						log.Error(set.Pattern.Pos.Line, set.Pattern.Pos.Column, "set statement variable '%s' not found in matches for production '%s'", varItem, production.Name)
					}
					err = CompileError{}
				}
			}
		}
	}

	if set.Pattern == nil && set.Value == nil {
		// should not be possible to get here since the parser should pick this up
		log.Error(set.Pos.Line, set.Pos.Column, "set statement is missing value (set to what?) in production '%s'", production.Name)
		err = CompileError{}
	}

	return
}

// validateRecallStatement checks a "recall" statement to verify the memory name.
func validateRecallStatement(recall *recallStatement, model *actr.Model, log *issues.Log, production *actr.Production) (err error) {
	vars := varsFromPattern(recall.Pattern)

	for _, varName := range vars {
		if varName == "?" {
			continue
		}

		match := production.LookupMatchByVariable(varName)
		if match == nil {
			log.Error(recall.Pos.Line, recall.Pos.Column, "recall statement variable '%s' not found in matches for production '%s'", varName, production.Name)
			err = CompileError{}
		}
	}

	return
}

// validateClearStatement checks a "clear" statement to verify the buffer names.
func validateClearStatement(clear *clearStatement, model *actr.Model, log *issues.Log, production *actr.Production) (err error) {
	bufferNames := clear.BufferNames

	for _, name := range bufferNames {
		buffer := model.LookupBuffer(name)
		if buffer == nil {
			log.Error(clear.Pos.Line, clear.Pos.Column, "buffer '%s' not found in production '%s'", name, production.Name)

			err = CompileError{}
			continue
		}
	}

	return
}

// validatePrintStatement is a placeholder for checking a "print" statement. Currently there are no checks.
func validatePrintStatement(print *printStatement, model *actr.Model, log *issues.Log, production *actr.Production) (err error) {
	if print.Args != nil {
		for _, v := range print.Args {
			if v.ID != nil {
				log.Error(print.Pos.Line, print.Pos.Column, "cannot use ID '%s' in print statement", *v.ID)
			} else if v.Var != nil {
				varItem := *v.Var
				match := production.LookupMatchByVariable(varItem)
				if match == nil {
					if varItem == "?" {
						log.Error(print.Pos.Line, print.Pos.Column, "cannot print anonymous var ('?') in production '%s'", production.Name)
					} else {
						log.Error(print.Pos.Line, print.Pos.Column, "print statement variable '%s' not found in matches for production '%s'", varItem, production.Name)
					}
					err = CompileError{}
				}
			}
		}
	}

	return
}

// validateVariableUsage verifies variable usage by counting how many times they are referenced.
func validateVariableUsage(log *issues.Log, match *match, do *do) {
	type ref struct {
		firstLine int // keep track of the first case of this variable for our output
		count     int // ref count
	}
	varRefCount := map[string]*ref{}

	// Walks a pattern to add all vars within
	addPatternRefs := func(p *pattern, insertIfNotFound bool) {
		vars := varsFromPattern(p)

		for _, varName := range vars {
			if varName == "?" {
				continue
			}

			if r, ok := varRefCount[varName]; ok {
				r.count++
			} else if insertIfNotFound {
				varRefCount[varName] = &ref{
					firstLine: p.Pos.Line,
					count:     1,
				}
			}
		}
	}

	// Walk the matches and store var ref counts
	for _, item := range match.Items {
		addPatternRefs(item.Pattern, true)
	}

	// Walk the do statements and add to var ref counts
	for _, statement := range *do.Statements {
		if statement.Set != nil {
			if statement.Set.Value != nil {
				if statement.Set.Value.Var != nil {
					varItem := *statement.Set.Value.Var
					if varItem != "?" {
						if r, ok := varRefCount[varItem]; ok {
							r.count++
						}
					}
				}
			} else { // pattern
				addPatternRefs(statement.Set.Pattern, false)
			}
		} else if statement.Recall != nil {
			addPatternRefs(statement.Recall.Pattern, false)
		} else if statement.Print != nil {
			for _, arg := range statement.Print.Args {
				if arg.Var != nil {
					varItem := *arg.Var
					if r, ok := varRefCount[varItem]; ok {
						r.count++
					}
				}
			}
		}
		// clear statement does not use variables
	}

	// Any var with only one reference should be anonymous ("?"), so add info to log
	for k, r := range varRefCount {
		if r.count == 1 {
			log.Error(r.firstLine, 0, "variable %s is not used - should be simplified to '?'", k)
		}
	}
}

// Get a slice of all the vars referenced in a pattern
func varsFromPattern(pattern *pattern) (vars []string) {
	for _, slot := range pattern.Slots {
		for _, slotItem := range slot.Items {
			if slotItem.Var != nil {
				vars = append(vars, *slotItem.Var)
			}
		}
	}

	return
}
