// Package version simply stores the build version string set using "-ldflags".
package version

// BuildVersion will be replaced during the build process using "-ldflags".
// See the Makefile.
var BuildVersion string = ""
