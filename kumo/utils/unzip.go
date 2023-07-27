package utils

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/samber/oops"
)

func Unzip(pathToZip, extractToPath string, bytesUnzipped chan<- int) (err error) {
	var (
		unzipGroup  = new(sync.WaitGroup)
		oopsBuilder = oops.Code("unzip_failed").
				With("pathToZip", pathToZip).
				With("extractToPath", extractToPath).
				With("bytesUnzipped", bytesUnzipped)

		reader  *zip.ReadCloser
		errChan chan error
		zipFile *zip.File
	)

	// Open the zip file and defer closing it
	if reader, err = zip.OpenReader(pathToZip); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to open zip file: %s", pathToZip)
		return
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Fatalf(
				"%+v",
				oopsBuilder.
					Wrapf(err, "failed to close zip reader: %#v", reader.File),
			)
		}
	}()

	// Declare error channel
	errChan = make(chan error, len(reader.File))

	// Unzip each file concurrently
	for _, zipFile = range reader.File {
		unzipGroup.Add(1)
		go func(zf *zip.File) {
			defer unzipGroup.Done()

			var (
				bytesCopied int64
			)

			if bytesCopied, err = unzipFile(zf, extractToPath); err != nil {
				err = oopsBuilder.
					With("bytesCopied", bytesCopied).
					With("zipFile", zf.Name).
					With("extractToPath", extractToPath).
					Wrapf(err, "failed to unzip file: %s", zf.Name)
				errChan <- err
				return
			}

			bytesUnzipped <- int(bytesCopied)
		}(zipFile)
	}

	// Wait for all files to be unzipped
	go func() {
		unzipGroup.Wait()
		close(errChan)
		close(bytesUnzipped)
	}()

	// Check for errors
	for err = range errChan {
		if err != nil {
			return
		}
	}

	return
}

func unzipFile(zf *zip.File, extractToPath string) (bytesCopied int64, err error) {
	var (
		oopsBuilder = oops.Code("unzipFile_failed").
				With("zipFile", zf).
				With("extractToPath", extractToPath)

		filePath        string
		destinationFile *os.File
		zippedFile      io.ReadCloser
	)

	// Check if file path is not vulnerable to Zip Slip
	filePath = filepath.Join(extractToPath, zf.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(extractToPath)+string(os.PathSeparator)) {
		err = oopsBuilder.
			Wrapf(err, "illegal file path: %s", filePath)
		return
	}

	// Check if file is a directory
	if zf.FileInfo().IsDir() {
		err = oopsBuilder.
			Errorf("is a directory: %s", zf.Name)
		return
	}

	// Create directory tree
	if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to create directory tree for: %s", filePath)
		return
	}

	// Create a destination file for unzipped content and defer closing it
	if destinationFile, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode()); err != nil {
		err = oopsBuilder.
			With("filePath", filePath).
			Wrapf(err, "failed to create destination file: %s", filePath)
		return
	}
	defer func() {
		if err := destinationFile.Close(); err != nil {
			log.Fatalf(
				"%+v",
				oopsBuilder.
					Wrapf(err, "failed to close destination file: %#v", destinationFile.Name()),
			)
		}
	}()

	// Unzip the content of a file and copy it to the destination file. Defer closing the zipped file
	if zippedFile, err = zf.Open(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to open zipped file: %s", zf.Name)
		return
	}
	defer func() {
		if err := zippedFile.Close(); err != nil {
			log.Fatalf(
				"%+v",
				oopsBuilder.
					Wrapf(err, "failed to close zipped file: %#v", zippedFile),
			)
		}
	}()

	if bytesCopied, err = io.Copy(destinationFile, zippedFile); err != nil {
		err = oopsBuilder.
			With("bytesCopied", bytesCopied).
			With("zippedFile", zippedFile).
			With("destinationFile", destinationFile).
			Wrapf(err, "failed to copy zipped file to destination file: %s", zf.Name)
		return
	}

	return
}
