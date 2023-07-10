package utils

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/host"
	"github.com/pkg/errors"
)

// Crafts the absolute path to a directory where the dependency will be downloaded
func CraftSingleZipPath(name string, s *host.Specs) (string, error) {
	cwd, err := GetCWD()
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		return "", err
	}
	depsDir := GetDepsDir()

	dzp := filepath.Join(cwd, depsDir, fmt.Sprintf("%s_%s_%s.zip", name, s.OS, s.ARCH))

	return dzp, nil
}
