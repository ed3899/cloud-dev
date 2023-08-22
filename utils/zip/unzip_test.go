package zip

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("unzipFile", func() {
	var (
		filename      = "test_file.txt"
		content       = "test content"
		mockZipReader *zip.Reader
		mockDestDir   string
	)

	BeforeEach(func() {
		// Temporary directory for testing
		cwd, err := os.Getwd()
		Expect(err).ToNot(HaveOccurred())
		mockDestDir = filepath.Join(cwd, "mock_dest")
	})

	AfterEach(func() {
		// Clean up the temporary directory after each test
		err := os.RemoveAll(mockDestDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("with a valid zip file", Label("unit"), func() {
		It("should unzip the file", func() {
			// Create a mock zip file in memory
			mockZipFile := &bytes.Buffer{}
			mockZipWriter := zip.NewWriter(mockZipFile)

			// Create a mock file
			mockDestFile, err := mockZipWriter.Create(filename)
			Expect(err).NotTo(HaveOccurred())

			// Write content to the mock file
			_, err = mockDestFile.Write([]byte(content))
			Expect(err).NotTo(HaveOccurred())

			err = mockZipWriter.Close()
			Expect(err).NotTo(HaveOccurred())

			// Create a mock zip reader
			mockZipReader, err = zip.NewReader(bytes.NewReader(mockZipFile.Bytes()), int64(mockZipFile.Len()))
			Expect(err).NotTo(HaveOccurred())

			// Unzip the mock file
			bytesUnzipped, err := unzipFile(mockZipReader.File[0], mockDestDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(bytesUnzipped).To(Equal(int64(len(content))))
		})
	})

	Context("with an invalid directory", Label("unit"), func() {
		It("should return an error", func() {
			// Create a mock zip file in memory
			mockZipFile := &bytes.Buffer{}
			mockZipWriter := zip.NewWriter(mockZipFile)

			// Create a mock file
			mockDestFile, err := mockZipWriter.Create(filename)
			Expect(err).NotTo(HaveOccurred())

			// Write content to the mock file
			_, err = mockDestFile.Write([]byte(content))
			Expect(err).NotTo(HaveOccurred())

			err = mockZipWriter.Close()
			Expect(err).NotTo(HaveOccurred())

			// Create a mock zip reader
			mockZipReader, err = zip.NewReader(bytes.NewReader(mockZipFile.Bytes()), int64(mockZipFile.Len()))
			Expect(err).NotTo(HaveOccurred())

			// Unzip the mock file
			_, err = unzipFile(mockZipReader.File[0], ":invalid_dir")
			Expect(err).To(HaveOccurred())
		})
	})
})
