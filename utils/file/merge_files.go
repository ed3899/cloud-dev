package file

import (
	"bufio"
	"fmt"
	"os"

	"github.com/samber/oops"
)

// Merges the content of the inputAbsFilePaths into a single file.
// The merged file is created at outputFileAbsPath.
func MergeFilesTo(
	outputFileAbsPath string,
	inputAbsFilePaths ...string,
) error {
	oopsBuilder := oops.
		Code("MergeFilesTo").
		In("utils").
		In("file").
		With("outputFileAbsPath", outputFileAbsPath).
		With("inputAbsFilePaths", inputAbsFilePaths)

	// Create a new file to write the merged content
	mergedFile, err := os.Create(outputFileAbsPath)
	if err != nil {
		err := oopsBuilder.Wrapf(err, "error creating merged file")

		return err
	}
	defer mergedFile.Close()

	for _, inputFilePath := range inputAbsFilePaths {
		// Open each file and append its content to the merged file
		inputFile, err := os.Open(inputFilePath)
		if err != nil {
			err := oopsBuilder.Wrapf(err, "error opening file %s", inputFilePath)

			return err
		}
		defer inputFile.Close()

		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			// Write each line to the merged file
			line := scanner.Text()
			_, err = fmt.Fprintf(mergedFile, "%s\n", line)
			if err != nil {
				err := oopsBuilder.Wrapf(err, "error writing line '%s' to merged file", line)

				return err
			}
		}

		err = scanner.Err()
		if err != nil {
			err := oopsBuilder.Wrapf(err, "error scanning file %s", inputFilePath)

			return err
		}
	}

	return nil
}
