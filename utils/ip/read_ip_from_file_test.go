package ip

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadIpFromFile", func() {

	Context("with a file containing an ip address", func() {
		var (
			testFileWithIp string
		)

		BeforeEach(func() {
			cwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			// Create a temporary file with an ip address
			fileContent := "Hello, this is an IP address: 192.168.1.1. Please read it."
			tmpFileWithIp, err := os.CreateTemp(cwd, "test_ip_file")
			Expect(err).NotTo(HaveOccurred())
			defer tmpFileWithIp.Close()

			_, err = tmpFileWithIp.Write([]byte(fileContent))
			Expect(err).NotTo(HaveOccurred())

			testFileWithIp = tmpFileWithIp.Name()
		})

		AfterEach(func() {
			// Clean up the temporary files
			err := os.Remove(testFileWithIp)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should read the IP address from a file", func() {
			ip, err := ReadIpFromFile(testFileWithIp)
			Expect(err).NotTo(HaveOccurred())
			Expect(ip).To(MatchRegexp("\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b"))
		})
	})

	Context("with a file not containing an ip address", func() {
		var (
			testFileWithoutIp string
		)

		BeforeEach(func() {
			cwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			// Create a temporary file without an IP address
			tmpFileWithoutIp, err := os.CreateTemp(cwd, "test_ip_file_no_ip")
			Expect(err).NotTo(HaveOccurred())
			defer tmpFileWithoutIp.Close()

			_, err = tmpFileWithoutIp.Write([]byte("No IP address here!"))
			Expect(err).NotTo(HaveOccurred())

			testFileWithoutIp = tmpFileWithoutIp.Name()
		})

		AfterEach(func() {
			err := os.Remove(testFileWithoutIp)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an error", func() {
			ip, err := ReadIpFromFile(testFileWithoutIp)
			Expect(err).To(HaveOccurred())
			Expect(ip).To(BeEmpty())
		})
	})

})
