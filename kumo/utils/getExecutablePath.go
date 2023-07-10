package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Returns the absolute path to the executable of the dependency.
// CWD is current working directory.
//
// - {CWD}/deps/{name}/{name}.exe
//
// - {CWD}/deps/{name}/{name}.exe
func GetExecutablePath(name string) (absExecutablePath string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		return "", err
	}

	return filepath.Join(cwd, GetDepsDir(), name, name+".exe"), nil
}
