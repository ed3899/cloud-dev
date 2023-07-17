package utils

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

func AttachCliToProcess(cmd *exec.Cmd) (err error) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

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

	errChan := make(chan error, 1)

	go func(src *io.ReadCloser, dest *os.File) {
		if _, err := io.Copy(dest, *src); err != nil {
			errx := errors.Wrap(err, "Error occurred while copying StdoutPipe to Stdout")
			errChan <- errx
			return
		}
	}(&stdout, os.Stdout)

	go func(src *io.ReadCloser, dest *os.File) {
		if _, err := io.Copy(dest, *src); err != nil {
			errx := errors.Wrap(err, "Error occurred while copying StderrPipe to Stderr")
			errChan <- errx
			return
		}
	}(&stderr, os.Stderr)

	go func() {
		for e := range errChan {
			if e != nil {
				errx := errors.Wrap(e, "Error occurred while copying std")
				log.Print(errx)
				log.Print("Sending interrupt signal to process...")
				TerminateCommand(cmd)
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case is := <-interrupt:
				if is != nil {
					log.Println("Interrupt signal received. Gracefully shutting down...")
					TerminateCommand(cmd)
					return
				}
			default:
				continue
			}
		}
	}()

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		err = errors.Wrap(err, "Error occurred while waiting for command to finish")
		close(errChan)
		close(interrupt)
		return err
	}

	close(errChan)
	close(interrupt)
	return nil
}
