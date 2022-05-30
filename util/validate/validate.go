// Package validate provides checks that fall between parsing & running.
package validate

import (
	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/issues"
)

// Goal adds a warning if we don't have a goal.
func Goal(model *actr.Model, initialGoal string, log *issues.Log) {
	if initialGoal == "" && !model.HasInitializer("buffer") {
		log.Warning(nil, "initial goal not provided and it was not initialized in the init section")
	}
}
