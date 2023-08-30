package tests

import (
	"os/exec"

	"github.com/ed3899/kumo/utils/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RunCmdAndStream", func() {
	var (
		testCmd *exec.Cmd
	)

	BeforeEach(func() {
		testCmd = exec.Command("python", "-c", "import time; print('Sleeping for 2 seconds...'); time.sleep(2); print('Done sleeping!')")
	})

	AfterEach(func() {
		testCmd = nil
	})

	It("should run a command and stream its output", Label("integration"), func() {
		err := cmd.RunCmdAndStream(testCmd)
		Expect(err).To(BeNil())
	})
})

var _ = Describe("TerminateCommand", func() {
	var (
		testCmd *exec.Cmd
	)

	BeforeEach(func() {
		testCmd = exec.Command("python", "-c", "import time; print('Sleeping for 2 seconds...'); time.sleep(2); print('Done sleeping!')")
	})

	AfterEach(func() {
		testCmd = nil
	})

	It("should terminate the command using SIGTERM or Kill", Label("integration"), func() {
		// Start the command execution
		err := testCmd.Start()
		Expect(err).To(BeNil())

		// Terminate the command
		cmd.TerminateCommand(testCmd)

		// Wait for the command to finish
		err = testCmd.Wait()
		Expect(err).To(HaveOccurred())
	})
})
