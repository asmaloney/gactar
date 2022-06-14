package framework

import (
	"errors"
	"fmt"

	"github.com/asmaloney/gactar/util/issues"
)

var (
	ErrModelMissingName = errors.New("model missing name")
)

type ErrBufferNotFound struct {
	BufferName string
	ModelName  string
}

func (e *ErrBufferNotFound) Error() string {
	return fmt.Sprintf("buffer %q not found in model %q", e.BufferName, e.ModelName)
}

type ErrExecutableNotSet struct {
	Name string
}

func (e *ErrExecutableNotSet) Error() string {
	return fmt.Sprintf("executable not set for %q", e.Name)
}

type ErrExecuteCommand struct {
	Output []byte
}

func (e *ErrExecuteCommand) Error() string {
	return fmt.Sprintf("execution failed:\n%s", string(e.Output))
}

type ErrModelGenerationFailed struct {
	Log *issues.Log
}

func (e *ErrModelGenerationFailed) Error() string {
	return e.Log.String()
}

type ErrModelValidationFailed struct {
	Log *issues.Log
}

func (e *ErrModelValidationFailed) Error() string {
	return e.Log.String()
}
