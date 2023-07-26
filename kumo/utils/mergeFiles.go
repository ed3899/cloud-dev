package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func MergeFilesTo(outputAbsFilePath string, inputAbsFilePaths ...string) (mergedFileAbsPath string, err error) {
	var (
		mergedFile    *os.File
		inputFilePath string
		inputFile     *os.File
		scanner       *bufio.Scanner
		line          string
	)

	// Create a new file to write the merged content
	if mergedFile, err = os.CreateTemp(filepath.Dir(outputAbsFilePath), filepath.Base(outputAbsFilePath)); err != nil {
		err = fmt.Errorf("error creating merged file: %v", err)
		return
	}
	defer func() {
		if errClosingMergedFile := mergedFile.Close(); err != nil {
			err = errors.Wrap(errClosingMergedFile, "error closing merged file")
		}
	}()

	for _, inputFilePath = range inputAbsFilePaths {
		// Open each file and append its content to the merged file
		if inputFile, err = os.Open(inputFilePath); err != nil {
			return "", errors.Wrapf(err, "error opening file %s", inputFilePath)
		}
		defer func() {
			if errClosingFile := inputFile.Close(); err != nil {
				err = errors.Wrap(errClosingFile, "error closing file")
			}
		}()

		scanner = bufio.NewScanner(inputFile)
		for scanner.Scan() {
			line = scanner.Text()
			if _, err = fmt.Fprintf(mergedFile, "%s\n", line); err != nil {
				return "", errors.Wrapf(err, "error writing line '%s' to merged file", line)
			}
		}

		if err = scanner.Err(); err != nil {
			return "", errors.Wrapf(err, "error scanning file %s", inputFilePath)
		}
	}

	mergedFileAbsPath = mergedFile.Name()

	return
}
