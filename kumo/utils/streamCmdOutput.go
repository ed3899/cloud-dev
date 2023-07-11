package utils

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/pkg/errors"
)

func RunCmdAndStreamOutput(cmd *exec.Cmd) (err error) {
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

	go func(dest *os.File, src *io.ReadCloser) {
		defer copyWg.Done()
		if _, err := io.Copy(dest, stdout); err != nil {
			err = errors.Wrap(err, "Error occurred while copying StdoutPipe to Stdout")
			errChan <- err
			return
		}
	}(os.Stdout, &stdout)

	go func(dest *os.File, src *io.ReadCloser) {
		defer copyWg.Done()
		if _, err := io.Copy(os.Stderr, stderr); err != nil {
			err = errors.Wrap(err, "Error occurred while copying StderrPipe to Stderr")
			errChan <- err
			return
		}
	}(os.Stderr, &stderr)

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
