package download

import (
	"os"

	"github.com/samber/oops"
)

func RemoveZip(
	download *Download,
) error {
	oopsBuilder := oops.
		Code("RemoveZip").
		With("download", download)

	err := os.Remove(download.Path.Zip)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", download.Path.Zip)

		return err
	}

	return nil
}
