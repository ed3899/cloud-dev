package utils

import (
	"bufio"
	"fmt"
	"os"
)

func MergeFiles(outputAbsFilePath string, inputAbsFilePaths ...string) (mergedFileAbsPath string, err error) {
	// Create a new file to write the merged content
	mergedFile, err := os.Create(outputAbsFilePath)
	if err != nil {
		err = fmt.Errorf("error creating merged file: %v", err)
		return
	}
	defer mergedFile.Close()

	// Open each file and append its content to the merged file
	for _, filePath := range inputAbsFilePaths {
		file, err := os.Open(filePath)
		if err != nil {
			err = fmt.Errorf("error opening file %s: %v", filePath, err)
			return "", err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			_, err := fmt.Fprintf(mergedFile, "%s\n", line)
			if err != nil {
				err = fmt.Errorf("error writing merged file: %v", err)
				return "", err
			}
		}

		if err := scanner.Err(); err != nil {
			err = fmt.Errorf("error scanning file %s: %v", filePath, err)
			return "", err
		}
	}

	mergedFileAbsPath = mergedFile.Name()

	return
}
