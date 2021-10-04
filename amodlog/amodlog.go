// Package amodlog provides info & error logging for parsing/compiling amod files.
package amodlog

import (
	"fmt"
	"io"
	"strings"
)

type level int

const (
	info level = iota
	err
)

type entry struct {
	level level
	line  int
	text  string
}

type Log struct {
	hasInfo  bool // does this log contain at least one info entry?
	hasError bool // does this log contain at least one error entry?
	entries  []entry
}

func (el *Log) addEntry(line int, l level, e string, a ...interface{}) {
	str := fmt.Sprintf(e, a...)
	el.entries = append(el.entries, entry{
		level: l,
		line:  line,
		text:  str,
	})
}

// New will create and return a new Log.
func New() *Log {
	return &Log{
		hasInfo:  false,
		hasError: false,
		entries:  []entry{},
	}
}

// HasInfo returns whether this log contains at least one info entry.
func (l *Log) HasInfo() bool {
	return l.hasInfo
}

// HasError returns whether this log contains at least one error entry.
func (l *Log) HasError() bool {
	return l.hasError
}

// Info will add a new info entry to the log.
func (l *Log) Info(line int, s string, a ...interface{}) {
	l.addEntry(line, info, s, a...)
	l.hasInfo = true
}

// Error will add a new error entry to the log.
func (l *Log) Error(line int, s string, a ...interface{}) {
	l.addEntry(line, err, s, a...)
	l.hasError = true
}

// String returns the log contents as a string. Each entry ends in a newline.
func (l *Log) String() string {
	b := new(strings.Builder)
	l.Write(b)
	return b.String()
}

// FirstEntry provides a way to get the text of the first log entry.
// This is used when parsing goals input by the user.
// For the UX we want to manage the output text differently.
// See amod.go ParseChunk()
func (l *Log) FirstEntry() string {
	if len(l.entries) == 0 {
		return "<INTERNAL ERROR: no Log entries exist>"
	}

	return l.entries[0].text
}

// Write will write the entire log. It will prepend INFO/ERROR and append
// line numbers (if any) to each log entry.
func (l *Log) Write(w io.Writer) {
	for _, entry := range l.entries {
		var str string

		switch entry.level {
		case info:
			str = "INFO: "
		case err:
			str = "ERROR: "
		}

		str += entry.text

		if entry.line != 0 {
			str += fmt.Sprintf(" (line %d)", entry.line)
		}

		str += "\n"

		w.Write([]byte(str))
	}
}
