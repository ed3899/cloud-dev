package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/samber/oops"
)

type IDownload interface {
	Path() string
	Name() string
}

func UnzipWith(
	zipOpenReader func(name string) (*zip.ReadCloser, error),
	unzipFile unzipFile,
) Unzip {
	oopsBuilder := oops.
		Code("UnzipWith")

	unzip := func(download IDownload, bytesUnzipped chan<- int) error {
		pathToZipDir := filepath.Dir(download.Path())
		pathToExecutable := filepath.Join(pathToZipDir, fmt.Sprintf("%s.exe", download.Name()))
		// Open the zip file and defer closing it
		reader, err := zip.OpenReader(pathToZip)
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to open zip file: %s", pathToZip)
			return err
		}
		return nil
	}

	return unzip
}

type Unzip func(download IDownload, bytesUnzipped chan<- int) error

func Unzip(
	pathToZip,
	extractToPath string,
	bytesUnzipped chan<- int,
) (
	err error,
) {
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
	defer func(reader *zip.ReadCloser) {
		if err := reader.Close(); err != nil {
			log.Fatalf(
				"%+v",
				oopsBuilder.
					Wrapf(err, "failed to close zip reader: %#v", reader.File),
			)
		}
	}(reader)

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

func unzipFileWith(
	osMkdirAll func(path string, perm os.FileMode) error,
	osOpenFile func(name string, flag int, perm os.FileMode) (*os.File, error),
	ioCopy func(dst io.Writer, src io.Reader) (written int64, err error),
	stringsHasPrefix func(s, prefix string) bool,
) unzipFile {
	oopsBuilder := oops.
		Code("unzipFileWith").
		In("utils").
		In("zip")

	unzipFile := func(zf *zip.File, extractToPath string) (int64, error) {
		// Check if file path is not vulnerable to Zip Slip
		filePath := filepath.Join(extractToPath, zf.Name)
		if !stringsHasPrefix(filePath, filepath.Clean(extractToPath)+string(os.PathSeparator)) {
			err := oopsBuilder.
				Errorf("illegal file path: %s", filePath)
			return -1, err
		}

		// Check if file is a directory
		if zf.FileInfo().IsDir() {
			err := oopsBuilder.
				Errorf("is a directory: %s", zf.Name)
			return -1, err
		}

		// Create directory tree
		err := osMkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to create directory tree for: %s", filePath)
			return -1, err
		}

		// Create a destination file for unzipped content and defer closing it
		destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode())
		if err != nil {
			err := oopsBuilder.
				With("filePath", filePath).
				Wrapf(err, "failed to create destination file: %s", filePath)
			return -1, err
		}
		defer destinationFile.Close()

		// Unzip the content of a file and copy it to the destination file. Defer closing the zipped file
		zippedFile, err := zf.Open()
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to open zipped file: %s", zf.Name)
			return -1, err
		}
		defer zippedFile.Close()

		bytesCopied, err := io.Copy(destinationFile, zippedFile)
		if err != nil {
			err := oopsBuilder.
				With("bytesCopied", bytesCopied).
				With("zippedFile", zippedFile).
				With("destinationFile", destinationFile).
				Wrapf(err, "failed to copy zipped file to destination file: %s", zf.Name)
			return -1, err
		}

		return bytesCopied, nil
	}

	return unzipFile
}

type unzipFile func(zf *zip.File, extractToPath string) (bytesCopied int64, err error)
