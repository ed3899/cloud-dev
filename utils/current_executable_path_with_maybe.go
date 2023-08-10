package utils

import (
	"github.com/samber/mo"
	"github.com/samber/oops"
)

func CurrentExecutablePath[
	CurrentExecutablePath ~string,
	OsExecutable ~func() (CurrentExecutablePath, error),
](
	osExecutable OsExecutable,
) mo.Result[CurrentExecutablePath] {
	oopsBuilder := oops.
		Code("CurrentExecutablePath")

	currentExecutablePath, err := osExecutable()
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to get current executable path")

		return mo.Err[CurrentExecutablePath](err)
	}

	return mo.Ok[CurrentExecutablePath](currentExecutablePath)
}
