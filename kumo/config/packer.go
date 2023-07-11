package config

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

func SetPackerPluginPath(absPath string) {
	err := os.Setenv("PACKER_PLUGIN_PATH", absPath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while setting PACKER_PLUGIN_PATH environment variable")
		log.Fatal(err)
	}
}
