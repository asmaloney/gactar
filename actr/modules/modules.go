// Package modules implements several ACT-R modules.
package modules

import (
	"github.com/asmaloney/gactar/actr/buffer"
	"github.com/asmaloney/gactar/actr/params"
)

// ModuleInterface provides an interface for the ACT-R concept of a "module".
type ModuleInterface interface {
	buffer.BufferInterface

	ModuleName() string

	SetParam(param *params.Param) (err params.ParamError)
}
