package framework

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type WriterHelper struct {
	File      *os.File
	Contents  *bytes.Buffer
	TabWriter *tabwriter.Writer
}

// KeyValueList is used to format output nicely with tabs using tabwriter.
type KeyValueList struct {
	list []keyValue
}

type keyValue struct {
	key   string
	value string
}

func (w *WriterHelper) InitWriterHelper(outputFileName string) (err error) {
	w.File, err = os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0770)
	if err != nil {
		return
	}

	w.Contents = new(bytes.Buffer)
	w.TabWriter = tabwriter.NewWriter(w.Contents, 0, 4, 1, '\t', 0)

	return
}

func (w *WriterHelper) CloseWriterHelper() (err error) {
	_, writeErr := w.File.Write(w.Contents.Bytes())

	// If we have a write error, we still want to try to Close().

	closeErr := w.File.Close()
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

	w.File = nil
	w.TabWriter = nil

	return
}

func (w WriterHelper) GetContents() []byte {
	return w.Contents.Bytes()
}

func (w WriterHelper) Write(e string, a ...interface{}) {
	str := fmt.Sprintf(e, a...)
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
