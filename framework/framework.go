package framework

import "gitlab.com/asmaloney/gactar/actr"

type Framework interface {
	Initialize() (err error)
	SetModel(model *actr.Model) (err error)

	Run(initialGoal string) (output []byte, err error)
	WriteModel(path string) (outputFileName string, err error)
}
