// Package modules implements several ACT-R modules.
package modules

import "github.com/asmaloney/gactar/actr/buffer"

// ModuleInterface provides an interface for the ACT-R concept of a "module".
type ModuleInterface interface {
	buffer.BufferInterface

	GetModuleName() string
}
