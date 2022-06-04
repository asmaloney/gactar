// Package issues provides info & error logging for parsing/compiling amod files.
package issues

import (
	"fmt"
	"io"
	"strings"
)

type level string

const (
	info    level = "info"
	warning level = "warning"
	err     level = "error"
)

type Location struct {
	Line        int `json:"line"`
	ColumnStart int `json:"columnStart"`
	ColumnEnd   int `json:"columnEnd"`
}

type Issue struct {
	Level level  `json:"level"`
	Text  string `json:"text"`

	*Location `json:"location"`
}

type IssueList = []Issue

type Log struct {
	hasError bool // does this log contain at least one error entry?
	issues   []Issue
}

// New will create and return a new Log.
func New() *Log {
	return &Log{
		hasError: false,
		issues:   []Issue{},
	}
}

// AllIssues returns a slice of all the current issues in the log.
func (l Log) AllIssues() IssueList {
	return l.issues
}

// HasIssues returns whether this log contains at least one entry.
func (l Log) HasIssues() bool {
	return len(l.issues) > 0
}

// HasError returns whether this log contains at least one error entry.
func (l Log) HasError() bool {
	return l.hasError
}

// Info will add a new info entry to the log.
func (l *Log) Info(location *Location, s string, a ...interface{}) {
	l.addEntry(location, info, s, a...)
}

// Warning will add a new info entry to the log.
func (l *Log) Warning(location *Location, s string, a ...interface{}) {
	l.addEntry(location, warning, s, a...)
}

// Error will add a new error entry to the log.
func (l *Log) Error(location *Location, s string, a ...interface{}) {
	l.addEntry(location, err, s, a...)
	l.hasError = true
}

// String returns the log contents as a string. Each entry ends in a newline.
func (l Log) String() string {
	b := new(strings.Builder)
	err := l.Write(b)
	if err != nil {
		return fmt.Sprintf("(could not write log to string: %s)", err.Error())
	}
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
func (l Log) Write(w io.Writer) error {
	for _, entry := range l.issues {
		var str string

		switch entry.Level {
		case info:
			str = "INFO: "
		case warning:
			str = "WARN: "
		case err:
			str = "ERROR: "
		}

		str += entry.Text

		if entry.Location != nil {
			str += fmt.Sprintf(" (line %d, col %d)", entry.Line, entry.ColumnStart)
		}

		str += "\n"

		_, err := w.Write([]byte(str))
		if err != nil {
			return err
		}
	}

	return nil
}

func (el *Log) addEntry(location *Location, l level, e string, a ...interface{}) {
	// If location is actually not set to anything, don't include it
	if location != nil && (*location == Location{}) {
		location = nil
	}
	str := fmt.Sprintf(e, a...)
	el.issues = append(el.issues, Issue{
		Level:    l,
		Text:     str,
		Location: location,
	})
}
