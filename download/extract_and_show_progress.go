package download

import (
	"github.com/samber/oops"
)

type IZipAndExecutable interface {
	Zip() string
	Executable() string
}

type IPathAndName interface {
	Path() IZipAndExecutable
	Name() string
}

func ExtractAndShowProgressWith(
	bar IIncrBy,
	utilsZipGetZipSize func(absPathToZip string) (size int64, err error),
	utilsZipUnzip func(download IPathAndName, bytesUnzipped chan<- int) error,
) {
	oopsBuilder := oops.
		Code("ExtractAndShowProgress")

	extractAndShowProgress := func(download IDownload) (Download, error) {
		downloadClone := download.Clone()

		zipSize, err := utilsZipGetZipSize(downloadClone.Path.zip)
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "failed to get zip size for: %v", dae.Download.AbsPath)
			return
		}
	}
}

type ExtractAndShowProgress func(IDownload) (Download, error)
