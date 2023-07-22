// Package decompress implements routines for unzipping and "untarring".
package decompress

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ErrZipInvalidFilePath struct {
	FilePath string
}

func (e ErrZipInvalidFilePath) Error() string {
	return fmt.Sprintf("invalid file path: %q", e.FilePath)
}

// Untar a file to a target directory.
func UntarFile(filePath, targetDir string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return Untar(file, targetDir)
}

// Untar a reader to a target directory.
func Untar(reader io.Reader, targetDir string) error {
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		filePath := header.Name
		if targetDir != "" {
			sanitized, pathErr := sanitizeExtractPath(targetDir, filePath)
			if pathErr != nil {
				return pathErr
			}
			filePath = sanitized
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(filePath); err != nil {
				if err := os.MkdirAll(filePath, 0755); err != nil {
					return err
				}
			}

		case tar.TypeReg:
			dstFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			for {
				_, fileErr := io.CopyN(dstFile, tr, 1024)
				if fileErr != nil {
					if errors.Is(fileErr, io.EOF) {
						break
					}
					return fileErr
				}
			}

			dstFile.Close()
		}
	}
}

// Unzip a file to a target directory.
func Unzip(filePath, targetDir string) (err error) {
	archive, err := zip.OpenReader(filePath)
	if err != nil {
		return
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := f.Name
		if targetDir != "" {
			sanitized, pathErr := sanitizeExtractPath(targetDir, filePath)
			if pathErr != nil {
				return pathErr
			}
			filePath = sanitized
		}

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return
			}
			continue
		}

		err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			return
		}

		dstFile, fileErr := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if fileErr != nil {
			return fileErr
		}

		fileInArchive, fileErr := f.Open()
		if fileErr != nil {
			return fileErr
		}

		for {
			_, fileErr := io.CopyN(dstFile, fileInArchive, 1024)
			if fileErr != nil {
				if errors.Is(fileErr, io.EOF) {
					break
				}
				return fileErr
			}
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return
}

// Sanitize archive file path to avoid "G305: Zip Slip vulnerability"
// See: https://snyk.io/research/zip-slip-vulnerability
func sanitizeExtractPath(targetDir, filePath string) (sanitized string, err error) {
	destPath := filepath.Join(targetDir, filePath)
	if !strings.HasPrefix(destPath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
		return "", &ErrZipInvalidFilePath{FilePath: filePath}
	}
	return destPath, nil
}
