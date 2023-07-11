package utils

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/pkg/errors"
)

func RunCmdAndStreamOutput(cmd *exec.Cmd) (err error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Stdin stat")
		return err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting StdinPipe")
		return err
	}

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

	if err := cmd.Start(); err != nil {
		err = errors.Wrap(err, "Error occurred while starting command")
		return err
	}

	copyWg := &sync.WaitGroup{}
	copyWg.Add(2)

	errChan := make(chan error)

	go func() {
		copyWg.Wait()
		close(errChan)
	}()

	// Check if there is any input available from stdin
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// If there is, copy it to stdinPipe
		copyWg.Add(1)

		go func(src *os.File, dest *io.WriteCloser) {
			defer copyWg.Done()

			if _, err := io.Copy(*dest, src); err != nil {
				err = errors.Wrap(err, "Error occurred while copying Stdin to StdinPipe")
				errChan <- err
				return
			}
			(*dest).Close()

		}(os.Stdin, &stdin)
	}

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

	go func() {
		if err := cmd.Wait(); err != nil {
			err = errors.Wrap(err, "Error occurred while waiting for command to finish")
			errChan <- err
			return
		}
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
