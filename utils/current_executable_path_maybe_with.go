package utils

import "github.com/samber/oops"

func CurrentExecutablePathMaybeWith[
	OsExecutable ~func() (string, error),
](
	osExecutable OsExecutable,
) (
	CurrentExecutablePath,
	error,
) {
	oopsBuilder := oops.
		Code("ExecutablePathMaybe")

	currentExecutablePath, err := osExecutable()
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to get current executable path")

		return nil, err
	}

	return func() string {
		return currentExecutablePath
	}, nil
}

type CurrentExecutablePath func() string
