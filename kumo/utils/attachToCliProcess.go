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

func AttachCliToProcess(cmd *exec.Cmd) (err error) {
	// Get StdoutPipe and StderrPipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting StdoutPipe")
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting StderrPipe")
		return err
	}

	// Start command
	if err := cmd.Start(); err != nil {
		err = errors.Wrap(err, "Error occurred while starting command")
		return err
	}

	// The wait group for all cmd related goroutines
	var mcg sync.WaitGroup
	// The error channel for all cmd related goroutines
	cmdErrChan := make(chan error, 1)
	// The done channel for all cmd related goroutines
	cmdGoRoutinesDone := make(chan bool, 1)

	// Copy command StdoutPipe to our Stdout
	mcg.Add(1)
	go func(src *io.ReadCloser, dest *os.File) {
		defer mcg.Done()
		if _, err := io.Copy(dest, *src); err != nil {
			errx := errors.Wrap(err, "Error occurred while copying StdoutPipe to Stdout")
			cmdErrChan <- errx
			return
		}
	}(&stdout, os.Stdout)

	// Copy command StderrPipe to our Stderr
	mcg.Add(1)
	go func(src *io.ReadCloser, dest *os.File) {
		defer mcg.Done()
		if _, err := io.Copy(dest, *src); err != nil {
			errx := errors.Wrap(err, "Error occurred while copying StderrPipe to Stderr")
			cmdErrChan <- errx
			return
		}
	}(&stderr, os.Stderr)

	// Wait for command to finish
	mcg.Add(1)
	go func() {
		defer mcg.Done()
		if err := cmd.Wait(); err != nil {
			err = errors.Wrap(err, "Error occurred while waiting for command to finish")
			cmdErrChan <- err
			cmdGoRoutinesDone <- false
			return
		}
	}()

	// Wait for all cmd related goroutines to finish
	go func() {
		defer close(cmdGoRoutinesDone)
		defer close(cmdErrChan)
		mcg.Wait()
		cmdGoRoutinesDone <- true
	}()

	// The aggregator wait group
	var sg sync.WaitGroup
	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	var mainErrChan = make(chan error, 1)

	sg.Add(1)
	go func() {
		defer sg.Done()
		defer close(mainErrChan)
		for {
			select {
			case err := <-cmdErrChan:
				if err != nil {
					err = errors.Wrap(err, "Error occurred while copying std")
					mainErrChan <- err
					TerminateCommand(cmd)
					return
				}

			case d := <-cmdGoRoutinesDone:
				if d {
					return
				}

			case s := <-signalChan:
				if s != nil {
					log.Println("You pressed Ctrl+C. Exiting...")
					TerminateCommand(cmd)
					return
				}
			default:
				continue
			}
		}
	}()

	sg.Wait()

	for err := range mainErrChan {
		if err != nil {
			return err
		}
	}

	return nil
}
