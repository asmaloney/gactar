package framework

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type WriterHelper struct {
	Contents  *bytes.Buffer
	TabWriter *tabwriter.Writer

	lineLen   int
	posInLine int
}

// KeyValueList is used to format output nicely with tabs using tabwriter.
type KeyValueList struct {
	list []keyValue
}

type keyValue struct {
	key   string
	value string
}

func (w *WriterHelper) InitWriterHelper() (err error) {
	w.Contents = new(bytes.Buffer)
	w.TabWriter = tabwriter.NewWriter(w.Contents, 0, 4, 1, '\t', 0)
	w.lineLen = 0
	w.posInLine = 0

	return
}

func (w *WriterHelper) WriteFile(outputFileName string) (err error) {
	file, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0770)
	if err != nil {
		return
	}

	_, writeErr := file.Write(w.Contents.Bytes())

	// If we have a write error, we still want to try to Close().

	closeErr := file.Close()
	if closeErr != nil {
		// If we also had an error on Close(), return both errors
		if writeErr != nil {
			err = fmt.Errorf("%s; %w", err.Error(), writeErr)
		} else {
			err = closeErr
		}
	} else {
		err = writeErr
	}

	return
}

func (w WriterHelper) GetContents() []byte {
	return w.Contents.Bytes()
}

// SetLineLen will cause newlines to be inserted at the line length.
// Note that this can be problematic for some things like comments or strings,
// so be careful how you use it.
func (w *WriterHelper) SetLineLen(lineLength int) {
	w.lineLen = lineLength
}

// ResetLineLen will reset the line length capability.
func (w *WriterHelper) ResetLineLen() {
	w.lineLen = 0
}

func (w *WriterHelper) Write(e string, a ...interface{}) {
	str := fmt.Sprintf(e, a...)

	if w.lineLen > 0 {
		strLen := len(str)
		w.posInLine += strLen
		if w.posInLine > w.lineLen {
			w.Contents.WriteString("\n")
			w.posInLine = strLen
		}
	}

	w.Contents.WriteString(str)
}

func (w WriterHelper) Writeln(e string, a ...interface{}) {
	w.Write(e+"\n", a...)
}

func (w WriterHelper) TabWrite(level int, list KeyValueList) {
	tabs := "\t"
	if level == 2 {
		tabs = "\t\t"
	} else if level > 2 {
		tabs = strings.Repeat("\t", level)
	}
	for _, item := range list.list {
		fmt.Fprintf(w.TabWriter, "%s%s\t%s\n", tabs, item.key, item.value)
	}

	w.TabWriter.Flush()
}

func (l *KeyValueList) Add(key, value string) {
	l.list = append(l.list, keyValue{
		key:   key,
		value: value,
	})
}
