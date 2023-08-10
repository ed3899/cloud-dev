package utils

import "github.com/samber/oops"

func CurrentExecutablePathMaybe(
	os_Executable func() (string, error),
) (
	func() string,
	error,
) {

	oopsBuilder := oops.
		Code("ExecutablePathMaybe")

	currentExecutablePath, err := os_Executable()
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to get current executable path")

		return nil, err
	}

	return func() string {
		return currentExecutablePath
	}, nil
}
