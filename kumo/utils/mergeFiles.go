package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/samber/oops"
)

func MergeFilesTo(outputAbsFilePath string, inputAbsFilePaths ...string) (mergedFileAbsPath string, err error) {
	var (
		oopsBuilder = oops.Code("merge_files_failed").
				With("outputAbsFilePath", outputAbsFilePath).
				With("inputAbsFilePaths", inputAbsFilePaths)

		mergedFile    *os.File
		inputFilePath string
		inputFile     *os.File
		scanner       *bufio.Scanner
		line          string
	)

	// Create a new file to write the merged content
	if mergedFile, err = os.CreateTemp(filepath.Dir(outputAbsFilePath), filepath.Base(outputAbsFilePath)); err != nil {
		err = oopsBuilder.Wrapf(err, "error creating merged file")
		return
	}
	defer func(mergedFile *os.File) {
		if err := mergedFile.Close(); err != nil {
			log.Fatalf(
				"%+v",
				oopsBuilder.
					Wrapf(err, "error closing merged file: %s", mergedFile.Name()),
			)
		}
	}(mergedFile)

	for _, inputFilePath = range inputAbsFilePaths {
		// Open each file and append its content to the merged file
		if inputFile, err = os.Open(inputFilePath); err != nil {
			err = oopsBuilder.Wrapf(err, "error opening file %s", inputFilePath)
			return
		}
		defer func(inputFile *os.File) {
			if err := inputFile.Close(); err != nil {
				log.Fatalf(
					"%+v",
					oopsBuilder.
						Wrapf(err, "error closing file: %s", inputFile.Name()),
				)
			}
		}(inputFile)

		scanner = bufio.NewScanner(inputFile)
		for scanner.Scan() {
			// Write each line to the merged file
			line = scanner.Text()
			if _, err = fmt.Fprintf(mergedFile, "%s\n", line); err != nil {
				err = oopsBuilder.Wrapf(err, "error writing line '%s' to merged file", line)
				return
			}
		}

		if err = scanner.Err(); err != nil {
			err = oopsBuilder.Wrapf(err, "error scanning file %s", inputFilePath)
			return
		}
	}

	mergedFileAbsPath = mergedFile.Name()

	return
}
