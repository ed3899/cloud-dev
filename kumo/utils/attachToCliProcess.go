package utils

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	// "os/signal"
	// "syscall"

	"github.com/pkg/errors"
)

func AttachCliToProcess(cmd *exec.Cmd) (err error) {
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

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

	var stop sync.WaitGroup
	stopChan := make(chan bool, 1)
	errChan := make(chan error, 1)

	stop.Add(1)
	go func() {
		defer stop.Done()
		reader := bufio.NewReader(os.Stdin)
		for {
			char, _, err := reader.ReadRune()
			if err != nil {
					log.Println("Error reading input:", err)
					return
			}

			switch char {
			case 'q':
					log.Println("You pressed 'q'. Exiting...")
					stopChan <- true
					return
			default:
					continue
			}
	}
	}()

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

	stop.Add(1)
	go func() {
		defer stop.Done()
		for {
			select {
			case s := <-stopChan:
				if s {
					log.Print("Gracefully shutting down...")
					TerminateCommand(cmd)
					return
				}
			case err := <-errChan:
				if err != nil {
					err = errors.Wrap(err, "Error occurred while copying std")
					log.Print(err)
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
		close(stopChan)
		// close(interrupt)
		return err
	}

	close(errChan)
	close(stopChan)
	// close(interrupt)
	return nil
}
