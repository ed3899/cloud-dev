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

func AttachToProcessStdAll(cmd *exec.Cmd) (err error) {
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

	copyWg := &sync.WaitGroup{}
	copyWg.Add(2)

	done := make(chan bool, 1)
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Close errChan and sigChan when copyWg is done
	go func() {
		copyWg.Wait()
		close(errChan)
		close(sigChan)
		close(done)
	}()

	// Copy StdoutPipe and StderrPipe to Stdout and Stderr
	go func(src *io.ReadCloser, dest *os.File) {
		defer copyWg.Done()
		if _, err := io.Copy(dest, *src); err != nil {
			err = errors.Wrap(err, "Error occurred while copying StdoutPipe to Stdout")
			errChan <- err
			return
		}
	}(&stdout, os.Stdout)

	go func(src *io.ReadCloser, dest *os.File) {
		defer copyWg.Done()
		if _, err := io.Copy(dest, *src); err != nil {
			err = errors.Wrap(err, "Error occurred while copying StderrPipe to Stderr")
			errChan <- err
			return
		}
	}(&stderr, os.Stderr)

	// Wait for command to finish
	go func() {
		if err := cmd.Wait(); err != nil {
			err = errors.Wrap(err, "Error occurred while waiting for command to finish")
			errChan <- err
			return
		}
		done <- true
	}()

// Wait for error, interrupt signal or done
OuterLoop:
	for {
		select {
		case errFromChan := <-errChan:
			switch {
			case errFromChan != nil:
				if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
					err = errors.Wrap(err, "Error occurred while sending interrupt signal to process")
					log.Fatal(err)
				}
				log.Fatal(err)
			default:
				continue OuterLoop
			}

		case <-sigChan:
			log.Println("Received interrupt signal, sending interrupt signal to process")
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				err = errors.Wrap(err, "Error occurred while sending interrupt signal to process")
				log.Fatal(err)
			}
			os.Exit(0)

		case <-done:
			return nil
		}
	}
}
