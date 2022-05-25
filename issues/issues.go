// Package issues provides info & error logging for parsing/compiling amod files.
package issues

import (
	"fmt"
	"io"
	"strings"
)

type level string

const (
	info level = "info"
	err  level = "error"
)

type Location struct {
	Line        int `json:"line"`
	ColumnStart int `json:"columnStart"`
	ColumnEnd   int `json:"columnEnd"`
}

type Issue struct {
	Level level  `json:"level"`
	Text  string `json:"text"`
	Location
}

type IssueList = []Issue

type Log struct {
	hasInfo  bool // does this log contain at least one info entry?
	hasError bool // does this log contain at least one error entry?
	issues   []Issue
}

// New will create and return a new Log.
func New() *Log {
	return &Log{
		hasInfo:  false,
		hasError: false,
		issues:   []Issue{},
	}
}

// AllIssues returns a slice of all the current issues in the log.
func (l Log) AllIssues() IssueList {
	return l.issues
}

// HasInfo returns whether this log contains at least one info entry.
func (l Log) HasInfo() bool {
	return l.hasInfo
}

// HasError returns whether this log contains at least one error entry.
func (l Log) HasError() bool {
	return l.hasError
}

// Info will add a new info entry to the log.
func (l *Log) Info(location *Location, s string, a ...interface{}) {
	l.addEntry(location, info, s, a...)
	l.hasInfo = true
}

// Error will add a new error entry to the log.
func (l *Log) Error(location *Location, s string, a ...interface{}) {
	l.addEntry(location, err, s, a...)
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
func (l Log) FirstEntry() string {
	if len(l.issues) == 0 {
		return "<INTERNAL ERROR: no Log entries exist>"
	}

	return l.issues[0].Text
}

// Write will write the entire log. It will prepend INFO/ERROR and append
// line numbers (if any) to each log entry.
func (l Log) Write(w io.Writer) {
	for _, entry := range l.issues {
		var str string

		switch entry.Level {
		case info:
			str = "INFO: "
		case err:
			str = "ERROR: "
		}

		str += entry.Text

		if entry.Line != 0 {
			str += fmt.Sprintf(" (line %d, col %d)", entry.Line, entry.ColumnStart)
		}

		str += "\n"

		w.Write([]byte(str))
	}
}

func (el *Log) addEntry(location *Location, l level, e string, a ...interface{}) {
	str := fmt.Sprintf(e, a...)
	el.issues = append(el.issues, Issue{
		Level:    l,
		Text:     str,
		Location: *location,
	})
}
