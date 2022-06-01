// Package validate provides checks that fall between parsing & running.
package validate

import (
	"github.com/asmaloney/gactar/actr"

	"github.com/asmaloney/gactar/util/issues"
)

// Goal adds a warning if we don't have a goal or adds info with the initial goal.
func Goal(model *actr.Model, initialGoal string, log *issues.Log) {
	initializer := model.LookupInitializer("goal")
	if initialGoal == "" && initializer == nil {
		log.Warning(nil, "initial goal not provided and it was not initialized in the init section")

		return
	}

	if initialGoal == "" {
		initialGoal = initializer.Pattern.String()
	}

	log.Info(nil, "initial goal is %s", initialGoal)
}
