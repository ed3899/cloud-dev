package utils

import (
	"os/exec"
	"sync"
	"testing"
)

func TestTerminateCommand(t *testing.T) {
	// Create a mock command for testing
	cmd := exec.Command("timeout", "10")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

	// The wait group for waiting cmd
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		cmd.Wait()
	}()

	// Terminate the command and be done with the go routine
	TerminateCommand(cmd)
	wg.Done()

	if cmd.ProcessState.ExitCode() != -1 {
		t.Errorf("Expected command to be terminated, but it is not")
	}
}
