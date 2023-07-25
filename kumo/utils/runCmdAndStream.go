package utils

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pkg/errors"
)

// Attaches the current CLI to the specified command.
// This allows the user to interact with the command as if it was run directly from the CLI.
//
// If the user presses **Ctrl+C**, this function signals the command is being attached to, to terminate.
// That way, the command will not keep running in the background.
//
// Example:
//
//	RunCmdAndStream(exec.Command("packer.exe", "build", "."))
//
//	`==> test.amazon-ebs.ubuntu: Prevalidating AMI Name: test`
//	`==> test.amazon-ebs.ubuntu: Found Image ID: ami-03f65b8614a860c29`
//	`==> test.amazon-ebs.ubuntu: Creating temporary keypair: packer_64b824bb-026f-af2c-184e-7097c138d520`
func RunCmdAndStream(cmd *exec.Cmd) (err error) {
	// TODO refactor
	var (
		cmdStdout io.ReadCloser
		cmdStderr io.ReadCloser

		cmdWg       = new(sync.WaitGroup)
		cmdErrChan  = make(chan error, 1)
		cmdDoneChan = make(chan bool, 1)

		aggregatorGroup = new(sync.WaitGroup)
		// The main error channel. This will be used to aggregate all errors from all cmd goroutines
		mainErrChan    = make(chan error, 1)
		signalChan     = make(chan os.Signal, 1)
		totalErr       error
		done           bool
		signalReceived os.Signal
	)
	// Get StdoutPipe
	cmdStdout, err = cmd.StdoutPipe()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting StdoutPipe")
		return err
	}
	defer cmdStdout.Close()

	// Get StderrPipe
	cmdStderr, err = cmd.StderrPipe()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting StderrPipe")
		return err
	}
	defer cmdStderr.Close()

	// Start command
	if err = cmd.Start(); err != nil {
		err = errors.Wrap(err, "Error occurred while starting command")
		return err
	}

	// Stream command StdoutPipe to our Stdout
	cmdWg.Add(1)
	go func(src *io.ReadCloser, dest *os.File) {
		defer cmdWg.Done()
		if _, err := io.Copy(dest, *src); err != nil {
			// In case of any streaming error, send the error to the error channel
			totalErr := errors.Wrap(err, "Error occurred while copying StdoutPipe to Stdout")
			cmdErrChan <- totalErr
			return
		}
	}(&cmdStdout, os.Stdout)

	// Stream command StderrPipe to our Stderr
	cmdWg.Add(1)
	go func(src *io.ReadCloser, dest *os.File) {
		defer cmdWg.Done()
		if _, err := io.Copy(dest, *src); err != nil {
			// In case of any streaming error, send the error to the error channel
			totalErr = errors.Wrap(err, "Error occurred while copying StderrPipe to Stderr")
			cmdErrChan <- totalErr
			return
		}
	}(&cmdStderr, os.Stderr)

	// Start a go routine to wait for the command to finish
	cmdWg.Add(1)
	go func() {
		defer cmdWg.Done()
		if err := cmd.Wait(); err != nil {
			// In case of any error while waiting for the command to finish, send the error to the error channel and send false to the done channel
			totalErr = errors.Wrap(err, "Error occurred while waiting for command to finish")
			cmdErrChan <- totalErr
			cmdDoneChan <- false
			return
		}
	}()

	// Start a go routine to wait for all cmd related goroutines to finish. When they finish, send true to the done channel. Now we can be sure no more errors will be sent to the cmd error channel
	go func() {
		defer close(cmdDoneChan)
		defer close(cmdErrChan)
		cmdWg.Wait()
		cmdDoneChan <- true
	}()

	// Notify the signal channel, to listen for Ctrl+C and other signals. This will allows us to terminate the command if the user presses Ctrl+C and pass that command termination to the cmd we are being attached to. This is important because if we don't terminate the command, it will keep running in the background.
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	aggregatorGroup.Add(1)
	go func() {
		defer aggregatorGroup.Done()
		defer close(mainErrChan)
		for {
			select {
			// If an error occurred while copying std, send it to the main error channel and terminate the command
			case err = <-cmdErrChan:
				if err != nil {
					err = errors.Wrap(err, "Error occurred while copying std")
					mainErrChan <- err
					TerminateCommand(cmd)
					return
				}

			// If the command finished successfully, return
			case done = <-cmdDoneChan:
				if done {
					return
				}

			// If the user pressed Ctrl+C or any other signal, terminate the command
			case signalReceived = <-signalChan:
				if signalReceived != nil {
					log.Println("You pressed Ctrl+C. Exiting...")
					TerminateCommand(cmd)
					return
				}
			default:
				continue
			}
		}
	}()

	aggregatorGroup.Wait()

	// Wait for all errors to be sent to the main error channel, if any.
	for err = range mainErrChan {
		if err != nil {
			return
		}
	}

	return
}

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
