package utils

import (
	"log"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
)

// Terminates the specified command. This command should not return an error.
// It is the end of the line.
func TerminateCommand(cmd *exec.Cmd) {
	// Attempt to send a SIGTERM signal to the process
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		totalErr := errors.Wrap(err, "Sending interrupt signal not supported")
		log.Print(totalErr)

		// Log that a kill signal will be sent instead
		log.Print("Sending kill signal instead")
		if err := cmd.Process.Kill(); err != nil {
			totalErr := errors.Wrap(err, "Sending kill signal not supported")
			log.Print(totalErr)
			return
		}

		return
	}
}