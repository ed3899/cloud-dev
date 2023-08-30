package tests

import (
	"os"

	"github.com/ed3899/kumo/utils/file"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsFilePresent", func() {
	Context("when a file exists", func() {
		var (
			_file *os.File
		)

		BeforeEach(func() {
			cwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			// Create a temporary file for testing
			_file, err = os.CreateTemp(cwd, "testfile.txt")
			defer _file.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			// Remove the temporary file
			GinkgoWriter.Println("Removing temporary file:", _file.Name())
			err := os.Remove(_file.Name())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return true", func() {
			Expect(file.IsFilePresent(_file.Name())).To(BeTrue())
		})
	})

	Context("when a file does not exist", func() {
		It("should return false", func() {
			Expect(file.IsFilePresent("/path/to/nonexistent/file")).To(BeFalse())
		})
	})
})
