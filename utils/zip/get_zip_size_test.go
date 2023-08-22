package zip

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetZipSize", func() {
	var (
		testZipFile *os.File
		content     string
	)

	BeforeEach(func() {
		// Create a temporary zip file for testing
		cwd, err := os.Getwd()
		Expect(err).ToNot(HaveOccurred())

		testZipFile, err = os.CreateTemp(cwd, "test.zip")
		Expect(err).ToNot(HaveOccurred())
		defer testZipFile.Close()

		// Write some data to the temporary zip file
		content = "Hello World!"
		_, err = testZipFile.WriteString(content)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Remove the temporary zip file
		err := os.Remove(testZipFile.Name())
		Expect(err).ToNot(HaveOccurred())
	})

	Context("with a valid zip path", Label("unit"), func() {
		It("should return the size of the zip file", func() {
			zipSize, err := GetZipSize(testZipFile.Name())
			Expect(err).ToNot(HaveOccurred())
			Expect(zipSize).To(Equal(int64(len(content))))
		})
	})

	Context("with an invalid zip path", Label("unit"), func() {
		It("should return an error", func() {
			_, err := GetZipSize("/invalid/path/to/zip")
			Expect(err).To(HaveOccurred())
		})
	})
})
