package examples

import (
	"embed"
)

// This works around the problem of embedding relative paths.
// Simply create a package for the examples and embed them!

//go:embed *.amod
var AMODExamples embed.FS
