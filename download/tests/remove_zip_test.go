package tests

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ed3899/kumo/download"
)

var _ = Describe("Download", func() {
	var (
		mockDownload *download.Download
	)

	Context("when the zip file exists", func() {
		BeforeEach(func() {
			cwd, err := os.Getwd()
			Expect(err).To(BeNil())

			file, err := os.CreateTemp(cwd, "file.zip")
			Expect(err).To(BeNil())
			defer file.Close()

			mockDownload = &download.Download{
				Path: &download.Path{
					Zip: file.Name(),
				},
			}
		})

		AfterEach(func() {
			os.RemoveAll(mockDownload.Path.Zip)
		})

		It("should remove the zip file", Label("unit"), func() {
			Expect(mockDownload.RemoveZip()).To(Succeed())

			// Expect the file to be removed
			_, err := os.Stat(mockDownload.Path.Zip)
			Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
		})
	})

	Context("when the zip file does not exist", func() {
		BeforeEach(func() {
			mockDownload = &download.Download{
				Path: &download.Path{
					Zip: "non_existent_file.zip",
				},
			}
		})

		It("should handle removal error", Label("unit"), func() {
			err := mockDownload.RemoveZip()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error occurred while removing"))
		})
	})
})
