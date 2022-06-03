// Package modules implements several ACT-R modules.
package modules

import "github.com/asmaloney/gactar/actr/buffer"

// Value mimics amod.fieldValue but without tokens.
type Value struct {
	ID     *string
	Str    *string
	Number *float64
}

type Param struct {
	Key   string
	Value Value
}

type ParamError = int

const (
	NoNumber ParamError = iota
	NumberRequired
	NumberMustBePositive

	UnrecognizedParam
)

// ModuleInterface provides an interface for the ACT-R concept of a "module".
type ModuleInterface interface {
	buffer.BufferInterface

	ModuleName() string

	SetParam(param *Param) (err ParamError)
}
