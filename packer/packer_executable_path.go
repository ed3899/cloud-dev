package packer

import (
	"fmt"

	"github.com/samber/oops"
)

func PackerExecutablePath(
	args PackerExecutableArgs,
) (string, error) {
	oopsBuilder := oops.Code("PackerExecutable")

	currentExecutablePath, err := args.os_Executable()
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to get current executable path")

		return "", err
	}

	return args.filepath_Join(
		currentExecutablePath,
		args.PackerName(),
		fmt.Sprintf("%s.exe", args.PackerName()),
	), nil
}

type PackerExecutableArgs struct {
	os_Executable func() (string, error)
	filepath_Join func(...string) string
	PackerName    func() string
}
