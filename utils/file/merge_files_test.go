package file

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MergeFilesTo", func() {
	var (
		content1 = "content of file 1\n"
		content2 = "content of file 2\n"

		mergedFilePath string
		file1          *os.File
		file2          *os.File
	)

	BeforeEach(func() {
		cwd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		mergedFilePath = filepath.Join(cwd, "output.txt")

		file1, err = os.CreateTemp(cwd, "testfile1.txt")
		defer file1.Close()
		Expect(err).NotTo(HaveOccurred())

		file2, err = os.CreateTemp(cwd, "testfile2.txt")
		defer file2.Close()
		Expect(err).NotTo(HaveOccurred())

		_, err = file1.WriteString(content1)
		Expect(err).NotTo(HaveOccurred())

		_, err = file2.WriteString(content2)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.Remove(file1.Name())
		Expect(err).NotTo(HaveOccurred())

		err = os.Remove(file2.Name())
		Expect(err).NotTo(HaveOccurred())

		err = os.Remove(mergedFilePath)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should merge input files into the output file", Label("unit"), func() {
		inputFilePaths := []string{
			file1.Name(),
			file2.Name(),
		}

		err := MergeFilesTo(mergedFilePath, inputFilePaths...)
		Expect(err).NotTo(HaveOccurred())

		// Check if the merged file was created correctly
		mergedContent, err := os.ReadFile(mergedFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(mergedContent)).To(Equal("content of file 1\ncontent of file 2\n"))
	})
})
