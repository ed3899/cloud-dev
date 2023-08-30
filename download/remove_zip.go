package download

import (
	"os"

	"github.com/samber/oops"
)

// Removes the associated zip file from the filesystem.
func (d *Download) RemoveZip() error {
	oopsBuilder := oops.
		In("download").
		Tags("Download").
		Code("RemoveZip")

	err := os.Remove(d.Path.Zip)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", d.Path.Zip)

		return err
	}

	return nil
}
