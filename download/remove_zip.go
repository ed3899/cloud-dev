package download

import (
	"os"

	"github.com/samber/oops"
)

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
