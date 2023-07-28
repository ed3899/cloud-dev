package cloud

import (
	"os"

	"github.com/samber/oops"
)

type RunDir struct {
	initial string
	target  string
}

func (rd *RunDir) GoTarget() (err error) {
	var (
		oopsBuilder = oops.
			Code("run_dir_go_target_failed")
	)

	if err = os.Chdir(rd.target); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while changing directory to %s", rd.target)
		return
	}

	return
}

func (rd *RunDir) GoInitial() (err error) {
	var (
		oopsBuilder = oops.
			Code("run_dir_go_initial_failed")
	)

	if err = os.Chdir(rd.initial); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while changing directory to %s", rd.initial)
		return
	}

	return
}
