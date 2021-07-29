package framework

import "gitlab.com/asmaloney/gactar/actr"

type Framework interface {
	SetModel(model *actr.Model) (err error)

	Run(initialGoal string) (err error)
	WriteModel() (outputFileName string, err error)
}
