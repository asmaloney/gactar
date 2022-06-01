package filesystem

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func DirExists(path string) bool {
	stat, err := os.Stat(path)
	return !os.IsNotExist(err) && stat.IsDir()
}

func CreateDir(path string) (err error) {
	err = os.MkdirAll(path, 0750)
	if err != nil && !os.IsExist(err) {
		return
	}

	return
}

// CheckForExecutable checks if an executable exists in the path.
func CheckForExecutable(exe string) (path string, err error) {
	path, err = exec.LookPath(exe)
	if err != nil {
		err = fmt.Errorf("cannot find '%s' in your path", exe)
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
