package utils

import (
	"os/exec"
	"sync"
	"time"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = Describe("TerminateCommand", func() {
	Context("TerminateCommand", func() {
		It("should terminate the command", func(ctx SpecContext) {
			// Create a mock command for testing
			cmd := exec.CommandContext(ctx, "timeout", "10")
			err := cmd.Start()
			Expect(err).To(BeNil())

			// The wait group for waiting cmd
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer GinkgoRecover()
				err = cmd.Wait()
				Expect(err).To(HaveOccurred())
			}()

			// Terminate the command
			TerminateCommand(cmd)
			// Wait for the command to finish. This should be instant because we terminated the command.
			wg.Wait()

			Expect(cmd.ProcessState.ExitCode()).To(Equal(1))
		}, SpecTimeout(time.Second * 10))
	})
})
