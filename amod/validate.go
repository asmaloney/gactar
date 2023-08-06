package amod

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/modules"

	"github.com/asmaloney/gactar/util/issues"
)

// varAndIndex is used to match var text to slot indices.
// This is used for issue locations.
type varAndIndex struct {
	text  string
	index int
}

// validateChunk checks the chunk name to ensure uniqueness and that it isn't using
// reserved names.
func validateChunk(model *actr.Model, log *issueLog, chunk *chunkDecl) (err error) {
	if actr.IsInternalChunkType(chunk.TypeName) {
		log.errorTR(chunk.Tokens, 1, 2, "cannot use reserved chunk type %q (chunks beginning with '_' are reserved)", chunk.TypeName)
		return ErrCompile
	}

	if actr.IsReservedType(chunk.TypeName) {
		log.errorTR(chunk.Tokens, 1, 2, "cannot use reserved chunk type %q", chunk.TypeName)
		return ErrCompile
	}

	c := model.LookupChunk(chunk.TypeName)
	if c != nil {
		log.errorTR(chunk.Tokens, 1, 2, "duplicate chunk type: '%s'", chunk.TypeName)
		return ErrCompile
	}

	return nil
}

func validateBufferInitialization(model *actr.Model, log *issueLog, moduleName string, buffer buffer.Interface, initializers []*namedInitializer) (err error) {
	for _, init := range initializers {
		pattern_err := validatePattern(model, log, init.Pattern)
		if pattern_err != nil {
			err = ErrCompile
			continue
		}
	}

	return
}

func validateModuleInitialization(model *actr.Model, log *issueLog, init *moduleInitializer) (err error) {
	moduleName := init.ModuleName
	module := model.LookupModule(moduleName)

	if module == nil {
		log.errorTR(init.Tokens, 0, 1, "module '%s' not found in initialization", moduleName)
		return ErrCompile
	}

	numBuffers := module.Buffers().Count()
	if numBuffers == 0 {
		log.errorTR(init.Tokens, 0, 1, "module '%s' does not have any buffers", moduleName)
		return ErrCompile
	}

	if len(init.InitPatterns) > 0 {
		if numBuffers > 1 {
			log.errorTR(init.Tokens, 0, 1, "module '%s' has more than one buffer - specify the buffer name", moduleName)
			return ErrCompile
		}

		buffer := module.Buffers().At(0)

		if !module.AllowsMultipleInit() && len(init.InitPatterns) > 1 {
			log.errorTR(init.InitPatterns[0].Tokens, 0, 1, "module %q should only have one pattern in initialization of buffer %q", moduleName, buffer.Name())
			return ErrCompile
		}

		err = validateBufferInitialization(model, log, moduleName, buffer, init.InitPatterns)
		if err != nil {
			err = ErrCompile
		}
	} else if len(init.BufferInitPatterns) > 0 {
		for _, bufferInit := range init.BufferInitPatterns {
			buff := model.LookupBuffer(bufferInit.BufferName)
			if buff == nil {
				log.errorTR(init.Tokens, 0, 1, "could not find buffer %q in module '%s' ", bufferInit.BufferName, moduleName)
				return ErrCompile
			}

			err = validateBufferInitialization(model, log, moduleName, buff, bufferInit.InitPatterns)
			if err != nil {
				err = ErrCompile
			}
		}
	}

	return
}

// validatePattern ensures that the pattern's chunk exists and that its number of slots match.
func validatePattern(model *actr.Model, log *issueLog, pattern *pattern) (err error) {
	chunkName := pattern.ChunkName
	chunk := model.LookupChunk(chunkName)
	if chunk == nil {
		log.errorTR(pattern.Tokens, 1, 2, "could not find chunk named '%s'", chunkName)
		return ErrCompile
	}

	if len(pattern.Slots) != chunk.NumSlots {
		s := "slots"
		if chunk.NumSlots == 1 {
			s = "slot"
		}
		log.errorT(pattern.Tokens, "invalid chunk - '%s' expects %d %s", chunkName, chunk.NumSlots, s)
		return ErrCompile
	}

	return
}

// validateBufferPatternMatch verifies several aspects of a chunk match item.
func validateBufferPatternMatch(item *matchBufferPatternItem, model *actr.Model, log *issueLog, production *actr.Production) (err error) {
	name := item.BufferName

	bufferInterface := model.LookupBuffer(name)
	if bufferInterface == nil {
		log.errorTR(item.Tokens, 0, 1, "buffer '%s' not found in production '%s'", name, production.Name)
		err = ErrCompile
		return
	}

	pattern := item.Pattern
	pattern_err := validatePattern(model, log, pattern)
	if pattern_err != nil {
		err = ErrCompile
	}

	// If we have constraints, check them
	if item.When != nil {
		for _, expr := range *item.When.Expressions {
			// Check that we haven't negated it in the pattern and then tried to constrain it further
			for _, slot := range pattern.Slots {
				if slot.Not && slot.Var != nil {
					if expr.LHS == *slot.Var {
						log.errorTR(expr.Tokens, 1, 2, "cannot further constrain a negated variable '%s'", expr.LHS)
						break
					}
				}
			}

			// Check that we aren't comparing to ourselves
			if expr.RHS.Var != nil && expr.LHS == *expr.RHS.Var {
				log.errorT(expr.RHS.Tokens, "cannot compare a variable to itself '%s'", expr.LHS)
			}
		}
	}

	return
}

// validateBufferStateMatch verifies a buffer match match item.
func validateBufferStateMatch(item *matchBufferStateItem, model *actr.Model, log *issueLog, production *actr.Production) (err error) {
	name := item.BufferName

	bufferInterface := model.LookupBuffer(name)
	if bufferInterface == nil {
		log.errorTR(item.Tokens, 0, 1, "buffer '%s' not found in production '%s'", name, production.Name)
		err = ErrCompile
	}

	if !buffer.IsValidState(item.State) {
		log.errorT(item.Tokens,
			"invalid state check '%s' for buffer '%s' in production '%s' (should be one of: %v)",
			item.State, name, production.Name, buffer.ValidStatesStr())
		err = ErrCompile
	}

	return
}

// validateModuleStateMatch verifies a module state match item.
func validateModuleStateMatch(item *matchModuleStateItem, model *actr.Model, log *issueLog, production *actr.Production) (err error) {
	name := item.ModuleName

	moduleInterface := model.LookupModule(name)
	if moduleInterface == nil {
		log.errorTR(item.Tokens, 0, 1, "module '%s' not found in production '%s'", name, production.Name)
		err = ErrCompile
	}

	if !modules.IsValidState(item.State) {
		log.errorT(item.Tokens,
			"invalid module state check '%s' for module '%s' in production '%s' (should be one of: %v)",
			item.State, name, production.Name, modules.ValidStatesStr())
		err = ErrCompile
	}

	return
}

// validateMatch verifies several aspects of a match item.
func validateMatch(match *match, model *actr.Model, log *issueLog, production *actr.Production) (err error) {
	if match == nil {
		return
	}

	for _, item := range match.Items {
		switch {
		case item.BufferPattern != nil:
			err = validateBufferPatternMatch(item.BufferPattern, model, log, production)

		case item.BufferState != nil:
			err = validateBufferStateMatch(item.BufferState, model, log, production)

		case item.ModuleState != nil:
			err = validateModuleStateMatch(item.ModuleState, model, log, production)
		}
	}

	return
}

// validateDo checks for multiple recall statements.
func validateDo(log *issueLog, production *production) {
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
		log.errorT([]lexer.Token{recallRef.token}, "only one recall statement per production is allowed in production '%s'", production.Name)
	}
}

// validateSetStatement checks a "set" statement to verify the buffer name & field indexing is correct.
// The production's matches have been constructed, so that's what we check against.
func validateSetStatement(set *setStatement, model *actr.Model, log *issueLog, production *actr.Production) (err error) {
	bufferName := set.BufferName
	buffer := model.LookupBuffer(bufferName)
	if buffer == nil {
		log.errorTR(set.Tokens, 1, 2, "buffer '%s' not found", bufferName)
		err = ErrCompile
	}

	if set.Slot != nil {
		// we have the form "set <buffer>.<slot name> to <value>"
		slotName := *set.Slot
		if set.Pattern != nil {
			log.errorTR(set.Tokens, 1, 3, "cannot set a slot ('%s.%s') to a pattern in production '%s'", bufferName, slotName, production.Name)
			err = ErrCompile
			return
		}

		match := production.LookupMatchByBuffer(bufferName)
		if match == nil {
			log.errorTR(set.Tokens, 1, 2, "match buffer '%s' not found in production '%s'", bufferName, production.Name)
			err = ErrCompile
			return
		}

		chunk := match.Pattern.Chunk
		if !chunk.HasSlot(slotName) {
			log.errorTR(set.Tokens, 3, 4, "slot '%s' does not exist in chunk type '%s' for match buffer '%s' in production '%s'", slotName, chunk.TypeName, bufferName, production.Name)
			err = ErrCompile
		}

		if set.Value.Var != nil {
			// Check set.Value.Var to ensure it exists
			varItem := *set.Value.Var
			match := production.LookupMatchByVariable(varItem)
			if match == nil {
				log.errorT(set.Value.Tokens, "set statement variable '%s' not found in matches for production '%s'", varItem, production.Name)
				err = ErrCompile
			}
		}
	} else {
		// we have the form "set <buffer> to <pattern>"
		if set.Value != nil {
			log.errorT(set.Value.Tokens, "buffer '%s' must be set to a pattern in production '%s'", bufferName, production.Name)
			err = ErrCompile
			return
		}

		chunkName := set.Pattern.ChunkName
		chunk := model.LookupChunk(chunkName)

		for slotIndex, slot := range set.Pattern.Slots {
			if slot.Var == nil {
				continue
			}

			if slot.Wildcard != nil {
				log.errorT(slot.Tokens, "cannot set '%s.%v' to wildcard ('*') in production '%s'", bufferName, chunk.SlotName(slotIndex), production.Name)
				err = ErrCompile
				continue
			}

			varItem := *slot.Var
			match := production.LookupMatchByVariable(varItem)
			if match == nil {
				log.errorT(slot.Tokens, "set statement variable '%s' not found in matches for production '%s'", varItem, production.Name)
				err = ErrCompile
			}
		}
	}

	return
}

// validateRecallStatement checks a "recall" statement to verify the memory name.
func validateRecallStatement(recall *recallStatement, model *actr.Model, log *issueLog, production *actr.Production) (err error) {
	pattern_err := validatePattern(model, log, recall.Pattern)
	if pattern_err != nil {
		err = ErrCompile
	}

	vars := varsFromPattern(recall.Pattern)

	for _, v := range vars {
		match := production.LookupMatchByVariable(v.text)
		if match == nil {
			log.errorT(recall.Pattern.Slots[v.index].Tokens, "recall statement variable '%s' not found in matches for production '%s'", v.text, production.Name)
			err = ErrCompile
		}
	}

	if recall.With != nil {
		buffer := model.Memory.BufferList.At(0)

		if buffer.RequestParameters() == nil {
			log.errorT(recall.With.Tokens, "recall 'with': buffer does not support any request parameters")
			err = ErrCompile
		} else {
			for _, param := range *recall.With.Expressions {
				key := param.Param

				if param.Value.Var != nil {
					log.errorT(param.Tokens, "recall 'with': parameter '%s'. Unexpected variable", key)
					err = ErrCompile
					continue
				}

				kv := argToKeyValue(key, param.Value)
				paramErr := buffer.RequestParameters().ValidateParam(kv)
				if paramErr != nil {
					log.errorT(param.Tokens,
						"recall 'with': %s.",
						paramErr.Error(),
					)
					err = ErrCompile
				}
			}
		}
	}

	return
}

// validateClearStatement checks a "clear" statement to verify the buffer names.
func validateClearStatement(clear *clearStatement, model *actr.Model, log *issueLog, production *actr.Production) (err error) {
	bufferNames := clear.BufferNames

	for _, name := range bufferNames {
		buffer := model.LookupBuffer(name)
		if buffer == nil {
			log.errorT(clear.Tokens, "buffer '%s' not found in production '%s'", name, production.Name)

			err = ErrCompile
			continue
		}
	}

	return
}

// validatePrintStatement checks a "print" statement's arguments.
//
//nolint:unparam // keeping the same function signature as the others
func validatePrintStatement(print *printStatement, model *actr.Model, log *issueLog, production *actr.Production) (err error) {
	if print.Args != nil {
		for _, arg := range print.Args {
			switch {
			case arg.Nil != nil:
				log.errorT(arg.Tokens, "cannot use nil in print statement")

			case arg.ID != nil:
				log.errorT(arg.Tokens, "cannot use ID '%s' in print statement", *arg.ID)

			case arg.Var != nil:
				varItem := *arg.Var
				match := production.LookupMatchByVariable(varItem)
				if match == nil {
					if varItem == "*" {
						log.errorT(arg.Tokens, "cannot print wildcard ('*') in production '%s'", production.Name)
					} else {
						log.errorT(arg.Tokens, "print statement variable '%s' not found in matches for production '%s'", varItem, production.Name)
					}
					err = ErrCompile
				}
			}
		}
	}

	return
}

// validateVariableUsage verifies variable usage by counting how many times they are referenced.
func validateVariableUsage(log *issueLog, match *match, do *do) {
	type ref struct {
		location *issues.Location // keep track of the first case of this variable for our output
		count    int              // ref count
	}
	varRefCount := map[string]*ref{}

	// Walks a pattern to add all vars within
	addPatternRefs := func(p *pattern, insertIfNotFound bool) {
		vars := varsFromPattern(p)

		for _, v := range vars {
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

	addWhenClauseRefs := func(w *whenExpression) {
		if r, ok := varRefCount[w.LHS]; ok {
			r.count++
		} else {
			log.errorTR(w.Tokens, 1, 2, "unknown variable %s in where clause", w.LHS)
			return
		}

		if w.RHS.Var != nil {
			if r, ok := varRefCount[*w.RHS.Var]; ok {
				r.count++
			} else {
				log.errorT(w.RHS.Tokens, "unknown variable %s in where clause", *w.RHS.Var)
				return
			}
		}
	}

	// Walk the matches and store var ref counts
	for _, match := range match.Items {
		// only need to consider chunk matchers
		if match.BufferPattern == nil {
			continue
		}

		addPatternRefs(match.BufferPattern.Pattern, true)

		if match.BufferPattern.When != nil {
			when := match.BufferPattern.When

			if when.Expressions != nil {
				for _, expr := range *when.Expressions {
					addWhenClauseRefs(expr)
				}
			}
		}
	}

	// Walk the do statements and add to var ref counts
	if do != nil {
		for _, statement := range *do.Statements {
			switch {
			case statement.Set != nil:
				if statement.Set.Value != nil {
					if statement.Set.Value.Var != nil {
						varItem := *statement.Set.Value.Var
						if r, ok := varRefCount[varItem]; ok {
							r.count++
						}
					}
				} else { // pattern
					addPatternRefs(statement.Set.Pattern, false)
				}

			case statement.Recall != nil:
				addPatternRefs(statement.Recall.Pattern, false)

			case statement.Print != nil:
				for _, arg := range statement.Print.Args {
					if arg.Var != nil {
						varItem := *arg.Var
						if r, ok := varRefCount[varItem]; ok {
							r.count++
						}
					}
				}

			default:
				// statement does not use variables
			}
		}
	}

	// Any var with only one reference should be wildcard ("*"), so add info to log
	for k, r := range varRefCount {
		if r.count == 1 {
			log.Error(r.location, "variable %s is not used - should be simplified to '*'", k)
		}
	}
}

// Get a slice of all the vars referenced in a pattern
func varsFromPattern(pattern *pattern) (vars []varAndIndex) {
	for i, slot := range pattern.Slots {
		if slot.Var != nil {
			vars = append(vars, varAndIndex{text: *slot.Var, index: i})
		}
	}

	return
}
