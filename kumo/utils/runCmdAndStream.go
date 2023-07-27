package utils

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/samber/oops"
	"go.uber.org/zap"
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
	var (
		oopsBuilder = oops.Code("run_cmd_and_stream_failed").
				With("cmd", cmd.Path)
		logger, _   = zap.NewProduction()
		cmdWg       = new(sync.WaitGroup)
		cmdErrChan  = make(chan error, 1)
		cmdDoneChan = make(chan bool, 1)

		aggregatorGroup = new(sync.WaitGroup)
		mainErrChan     = make(chan error, 1)
		signalChan      = make(chan os.Signal, 1)

		cmdStdout io.ReadCloser
		cmdStderr io.ReadCloser
	)

	// Zap logger setup
	defer logger.Sync()

	// Get StdoutPipe
	if cmdStdout, err = cmd.StdoutPipe(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while getting StdoutPipe for command '%s'", cmd.Path)
		return
	}

	// Get StderrPipe
	if cmdStderr, err = cmd.StderrPipe(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while getting StderrPipe for command '%s'", cmd.Path)
		return
	}

	// Start command
	if err = cmd.Start(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while starting command '%s'", cmd.Path)
		return
	}

	// Stream command StdoutPipe to our Stdout
	cmdWg.Add(1)
	go func(src *io.ReadCloser, dest *os.File) {
		defer cmdWg.Done()
		if _, err := io.Copy(dest, *src); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while copying StdoutPipe to Stdout for command '%s'", cmd.Path)
			cmdErrChan <- err
			return
		}
	}(&cmdStdout, os.Stdout)

	// Stream command StderrPipe to our Stderr
	cmdWg.Add(1)
	go func(src *io.ReadCloser, dest *os.File) {
		defer cmdWg.Done()
		if _, err := io.Copy(dest, *src); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while copying StderrPipe to Stderr for command '%s'", cmd.Path)
			cmdErrChan <- err
			return
		}
	}(&cmdStderr, os.Stderr)

	// Start a go routine to wait for the command to finish
	cmdWg.Add(1)
	go func() {
		defer cmdWg.Done()
		if err := cmd.Wait(); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while waiting for command '%s' to finish", cmd.Path)
			cmdErrChan <- err
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
			case err := <-cmdErrChan:
				if err != nil {
					err = oopsBuilder.
						Wrapf(err, "Error occurred while copying std for command '%s'", cmd.Path)
					mainErrChan <- err
					TerminateCommand(cmd)
					return
				}

			// If the command finished successfully, return
			case done := <-cmdDoneChan:
				if done {
					return
				}

			// If the user pressed Ctrl+C or any other signal, terminate the command
			case signalReceived := <-signalChan:
				if signalReceived != nil {
					logger.Info("Exiting because of signal received...", zap.String("signal", signalReceived.String()))
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

// Terminates the specified command. This command should not return an error and
// should not rely on external packages, this is in order to avoid additional error handling.
// It is the end of the line.
func TerminateCommand(cmd *exec.Cmd) {
	// Attempt to send a SIGTERM signal to the process
	if signalErr := cmd.Process.Signal(syscall.SIGTERM); signalErr != nil {
		log.Printf("Sending interrupt signal failed with error: %v\n", signalErr)
		log.Println("Sending kill signal instead")
		if killErr := cmd.Process.Kill(); killErr != nil {
			log.Printf("Sending kill signal failed with error: %v\n", killErr)
			log.Fatal("Failed to terminate command")
		}
	}
}
