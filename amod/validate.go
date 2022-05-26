package amod

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/issues"
)

// varAndIndex is used to match var text to slot indices.
// This is used for issue locations.
type varAndIndex struct {
	text  string
	index int
}

// validateChunk checks the chunk name to ensure uniqueness and that it isn't using
// reserved names.
func validateChunk(model *actr.Model, log *Log, chunk *chunkDecl) (err error) {
	if actr.IsInternalChunkName(chunk.Name) {
		log.ErrorTR(chunk.Tokens, 1, 2, "cannot use reserved chunk name '%s' (chunks beginning with '_' are reserved)", chunk.Name)
		return CompileError{}
	}

	c := model.LookupChunk(chunk.Name)
	if c != nil {
		log.ErrorTR(chunk.Tokens, 1, 2, "duplicate chunk name: '%s'", chunk.Name)
		return CompileError{}
	}

	return nil
}

func validateInitialization(model *actr.Model, log *Log, init *initialization) (err error) {
	name := init.Name
	module := model.LookupModule(name)

	if module == nil {
		log.ErrorTR(init.Tokens, 0, 1, "module '%s' not found in initialization", name)
		return CompileError{}
	}

	if !module.AllowsMultipleInit() && len(init.InitPatterns) > 1 {
		log.ErrorTR(init.Tokens, 0, 1, "module '%s' should only have one pattern in initialization", name)
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
func validatePattern(model *actr.Model, log *Log, pattern *pattern) (err error) {
	chunkName := pattern.ChunkName
	chunk := model.LookupChunk(chunkName)
	if chunk == nil {
		log.ErrorTR(pattern.Tokens, 1, 2, "could not find chunk named '%s'", chunkName)
		return CompileError{}
	}

	if len(pattern.Slots) != chunk.NumSlots {
		s := "slots"
		if chunk.NumSlots == 1 {
			s = "slot"
		}
		log.ErrorT(pattern.Tokens, "invalid chunk - '%s' expects %d %s", chunkName, chunk.NumSlots, s)
		return CompileError{}
	}

	return
}

// validateMatch verifies several aspects of a match item.
func validateMatch(match *match, model *actr.Model, log *Log, production *actr.Production) (err error) {
	if match == nil {
		return
	}

	for _, item := range match.Items {
		name := item.Name

		buffer := model.LookupBuffer(name)

		if buffer == nil {
			log.ErrorT(item.Tokens, "buffer '%s' not found in production '%s'", name, production.Name)
			err = CompileError{}
			continue
		}

		pattern := item.Pattern
		if pattern == nil {
			log.ErrorT(item.Tokens, "invalid pattern for '%s' in production '%s'", name, production.Name)
			err = CompileError{}
			continue
		}

		pattern_err := validatePattern(model, log, pattern)
		if pattern_err != nil {
			err = CompileError{}
		}

		// check _status chunks to ensure they have one of the allowed tests
		if pattern.ChunkName == "_status" {
			slot := *pattern.Slots[0].Items[0].ID

			if !actr.IsValidBufferState(slot) {
				log.ErrorT(item.Tokens, "invalid _status '%s' for '%s' in production '%s' (should be %v)", slot, name, production.Name, actr.ValidBufferStatesStr())
				err = CompileError{}
			}
		}
	}

	return
}

// validateDo checks for multiple recall statements.
func validateDo(log *Log, production *production) {
	type ref struct {
		token lexer.Token // keep track of the "recall" token from last case
		count int         // ref count
	}

	recallRef := ref{
		count: 0,
	}

	for _, statement := range *production.Do.Statements {
		if statement.Recall != nil {
			recallRef.token = statement.Tokens[0]
			recallRef.count++
		}
	}

	if recallRef.count > 1 {
		log.ErrorT([]lexer.Token{recallRef.token}, "only one recall statement per production is allowed in production '%s'", production.Name)
	}
}

// validateSetStatement checks a "set" statement to verify the buffer name & field indexing is correct.
// The production's matches have been constructed, so that's what we check against.
func validateSetStatement(set *setStatement, model *actr.Model, log *Log, production *actr.Production) (err error) {
	bufferName := set.BufferName
	buffer := model.LookupBuffer(bufferName)
	if buffer == nil {
		log.ErrorT(set.Tokens, "buffer '%s' not found in production '%s'", bufferName, production.Name)
		err = CompileError{}
	}

	if set.Slot != nil {
		// we have the form "set <buffer>.<slot name> to <value>"
		slotName := *set.Slot
		if set.Pattern != nil {
			log.ErrorT(set.Tokens, "cannot set a slot ('%s') to a pattern in match buffer '%s' in production '%s'", slotName, bufferName, production.Name)
			err = CompileError{}
		} else {
			match := production.LookupMatchByBuffer(bufferName)

			if match == nil {
				log.ErrorT(set.Tokens, "match buffer '%s' not found in production '%s'", bufferName, production.Name)
				err = CompileError{}
			} else {
				chunk := match.Pattern.Chunk
				if chunk == nil {
					log.ErrorT(set.Tokens, "chunk does not exist in match buffer '%s' in production '%s'", bufferName, production.Name)
					err = CompileError{}
				} else {
					if !chunk.HasSlot(slotName) {
						log.ErrorT(set.Tokens, "slot '%s' does not exist in chunk '%s' for match buffer '%s' in production '%s'", slotName, chunk.Name, bufferName, production.Name)
						err = CompileError{}
					}

					if set.Value.Var != nil {
						// Check set.Value.Var to ensure it exists
						varItem := *set.Value.Var
						match := production.LookupMatchByVariable(varItem)
						if match == nil {
							if varItem == "?" {
								log.ErrorT(set.Value.Tokens, "cannot set '%s.%s' to anonymous var ('?') in production '%s'", bufferName, slotName, production.Name)
							} else {
								log.ErrorT(set.Value.Tokens, "set statement variable '%s' not found in matches for production '%s'", varItem, production.Name)
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
			log.ErrorT(set.Tokens, "buffer '%s' must be set to a pattern in production '%s'", bufferName, production.Name)
			err = CompileError{}
		} else {
			chunkName := set.Pattern.ChunkName
			chunk := model.LookupChunk(chunkName)

			for slotIndex, slot := range set.Pattern.Slots {
				if len(slot.Items) > 1 {
					log.ErrorT(set.Pattern.Tokens, "cannot set '%s.%v' to compound var in production '%s'", bufferName, chunk.SlotName(slotIndex), production.Name)
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
						log.ErrorT(set.Pattern.Tokens, "cannot set '%s.%v' to anonymous var ('?') in production '%s'", bufferName, chunk.SlotName(slotIndex), production.Name)
					} else {
						log.ErrorT(set.Pattern.Tokens, "set statement variable '%s' not found in matches for production '%s'", varItem, production.Name)
					}
					err = CompileError{}
				}
			}
		}
	}

	if set.Pattern == nil && set.Value == nil {
		// should not be possible to get here since the parser should pick this up
		log.ErrorT(set.Tokens, "set statement is missing value (set to what?) in production '%s'", production.Name)
		err = CompileError{}
	}

	return
}

// validateRecallStatement checks a "recall" statement to verify the memory name.
func validateRecallStatement(recall *recallStatement, model *actr.Model, log *Log, production *actr.Production) (err error) {
	pattern_err := validatePattern(model, log, recall.Pattern)
	if pattern_err != nil {
		err = CompileError{}
	}

	vars := varsFromPattern(recall.Pattern)

	for _, v := range vars {
		if v.text == "?" {
			continue
		}

		match := production.LookupMatchByVariable(v.text)
		if match == nil {
			log.ErrorT(recall.Pattern.Slots[v.index].Tokens, "recall statement variable '%s' not found in matches for production '%s'", v.text, production.Name)
			err = CompileError{}
		}
	}

	return
}

// validateClearStatement checks a "clear" statement to verify the buffer names.
func validateClearStatement(clear *clearStatement, model *actr.Model, log *Log, production *actr.Production) (err error) {
	bufferNames := clear.BufferNames

	for _, name := range bufferNames {
		buffer := model.LookupBuffer(name)
		if buffer == nil {
			log.ErrorT(clear.Tokens, "buffer '%s' not found in production '%s'", name, production.Name)

			err = CompileError{}
			continue
		}
	}

	return
}

// validatePrintStatement is a placeholder for checking a "print" statement. Currently there are no checks.
func validatePrintStatement(print *printStatement, model *actr.Model, log *Log, production *actr.Production) (err error) {
	if print.Args != nil {
		for _, v := range print.Args {
			if v.ID != nil {
				log.ErrorT(print.Tokens, "cannot use ID '%s' in print statement", *v.ID)
			} else if v.Var != nil {
				varItem := *v.Var
				match := production.LookupMatchByVariable(varItem)
				if match == nil {
					if varItem == "?" {
						log.ErrorT(print.Tokens, "cannot print anonymous var ('?') in production '%s'", production.Name)
					} else {
						log.ErrorT(print.Tokens, "print statement variable '%s' not found in matches for production '%s'", varItem, production.Name)
					}
					err = CompileError{}
				}
			}
		}
	}

	return
}

// validateVariableUsage verifies variable usage by counting how many times they are referenced.
func validateVariableUsage(log *Log, match *match, do *do) {
	type ref struct {
		location *issues.Location // keep track of the first case of this variable for our output
		count    int              // ref count
	}
	varRefCount := map[string]*ref{}

	// Walks a pattern to add all vars within
	addPatternRefs := func(p *pattern, insertIfNotFound bool) {
		vars := varsFromPattern(p)

		for _, v := range vars {
			if v.text == "?" {
				continue
			}

			if r, ok := varRefCount[v.text]; ok {
				r.count++
			} else if insertIfNotFound {
				tokens := p.Slots[v.index].Tokens
				varRefCount[v.text] = &ref{
					location: tokensToLocation(tokens),
					count:    1,
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
			log.Error(r.location, "variable %s is not used - should be simplified to '?'", k)
		}
	}
}

// Get a slice of all the vars referenced in a pattern
func varsFromPattern(pattern *pattern) (vars []varAndIndex) {
	for i, slot := range pattern.Slots {
		for _, slotItem := range slot.Items {
			if slotItem.Var != nil {
				vars = append(vars, varAndIndex{text: *slotItem.Var, index: i})
			}
		}
	}

	return
}
