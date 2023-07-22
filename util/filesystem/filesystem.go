// Package filesystem implements some functions and error handling for working
// with files and directories.
package filesystem

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

type ErrDirDoesNotExist struct {
	DirName string
}

func (e ErrDirDoesNotExist) Error() string {
	return fmt.Sprintf("directory does not exist: %q", e.DirName)
}

type ErrFileDoesNotExist struct {
	FileName string
}

func (e ErrFileDoesNotExist) Error() string {
	return fmt.Sprintf("file does not exist: %q", e.FileName)
}

type ErrExeNotFound struct {
	ExeName string
	Path    string
}

func (e ErrExeNotFound) Error() string {
	return fmt.Sprintf("cannot find %q in your path:\n%q", e.ExeName, e.Path)
}

// DirExists returns true if the given path exists and is a directory.
func DirExists(path string) bool {
	stat, err := os.Stat(path)
	return !os.IsNotExist(err) && stat.IsDir()
}

// CreateDir creates a directory if it does not exist.
func CreateDir(path string) (err error) {
	err = os.MkdirAll(path, 0750)
	if err != nil && !os.IsExist(err) {
		return
	}

	return
}

// DownloadFile downloads a file from a URL.
func DownloadFile(url *url.URL, filePath string) (err error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return
}

// CheckForExecutable checks if an executable exists in the path.
func CheckForExecutable(exe string) (path string, err error) {
	path, err = exec.LookPath(exe)
	if err != nil {
		err = &ErrExeNotFound{ExeName: exe, Path: os.Getenv("PATH")}
		return "", err
	}

	return
}

// RemoveFile removes the given file if it exists.
func RemoveFile(filePath string) error {
	_, err := os.Stat(filePath)
	if err == nil {
		return os.Remove(filePath)
	}

	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}
