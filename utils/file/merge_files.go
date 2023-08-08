package file

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/samber/oops"
)

func MergeFilesTo(
	outputDirAbsPath string,
	inputAbsFilePaths ...string,
) (
	mergedFileName string,
	err error,
) {
	var (
		oopsBuilder = oops.Code("merge_files_failed").
				With("outputDirAbsPath", outputDirAbsPath).
				With("inputAbsFilePaths", inputAbsFilePaths)

		mergedFile    *os.File
		inputFilePath string
		inputFile     *os.File
		scanner       *bufio.Scanner
		line          string
	)

	// Create a new file to write the merged content
	if mergedFile, err = os.CreateTemp(outputDirAbsPath, "temp_file"); err != nil {
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

	mergedFileName = mergedFile.Name()

	return
}

type MergeFilesToF func(
	outputDirAbsPath string,
	inputAbsFilePaths ...string,
) (
	mergedFileName string,
	err error,
)
