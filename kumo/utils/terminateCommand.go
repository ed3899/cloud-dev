package utils

import (
	"log"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
)

func TerminateCommand(cmd *exec.Cmd) {
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		errx := errors.Wrap(err, "Sending interrupt signal not supported")
		log.Print(errx)

		log.Print("Sending kill signal instead")
		if err := cmd.Process.Kill(); err != nil {
			errx := errors.Wrap(err, "Sending kill signal not supported")
			log.Print(errx)
			return
		}

		return
	}
}