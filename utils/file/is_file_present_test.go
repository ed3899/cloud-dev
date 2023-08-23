package file

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsFilePresent", func() {
	Context("when a file exists", func() {
		var (
			file *os.File
		)

		BeforeEach(func() {
			cwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			// Create a temporary file for testing
			file, err = os.CreateTemp(cwd, "testfile.txt")
			defer file.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			// Remove the temporary file
			GinkgoWriter.Println("Removing temporary file:", file.Name())
			err := os.Remove(file.Name())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return true", func() {
			Expect(IsFilePresent(file.Name())).To(BeTrue())
		})
	})

	Context("when a file does not exist", func() {
		It("should return false", func() {
			Expect(IsFilePresent("/path/to/nonexistent/file")).To(BeFalse())
		})
	})
})
