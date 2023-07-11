package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/host"
	"github.com/pkg/errors"
)

// Crafts the absolute path to a directory where the dependency will be downloaded
func CraftSingleZipPath(name string, s *host.Specs) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		return "", err
	}
	depsDir := GetDependenciesDirName()

	dzp := filepath.Join(cwd, depsDir, fmt.Sprintf("%s_%s_%s.zip", name, s.OS, s.ARCH))

	return dzp, nil
}
