package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetPackerPluginPath() (ppp string) {
	cwd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		log.Fatal(err)
	}

	ppp = filepath.Join(cwd, "packer", "plugins")
	return ppp
}
