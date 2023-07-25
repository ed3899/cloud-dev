package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = Describe("CheckExistanceOfDir", func() {
	var tempDirPath string

	BeforeEach(func() {
		// Get current working directory
		cwd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		// Create a temporary directory for testing
		tempDirPath = filepath.Join(cwd, "tmp")
		err = os.Mkdir(tempDirPath, 0755)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.Remove(tempDirPath)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("TestDirExists", func() {
		When("the directory exists", func() {
			It("should return true", func() {
				exist, err := DirExist(tempDirPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(exist).To(BeTrue())
			})
		})
		When("the directory does not exist", func() {
			It("should return false", func() {
				exist, err := DirExist("path/to/nonexisting/dir")
				Expect(err).NotTo(HaveOccurred())
				Expect(exist).To(BeFalse())
			})
		})
	})

	Context("TestDirNotExists", func() {
		When("the directory exists", func() {
			It("should return false", func() {
				exist, err := DirNotExist(tempDirPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(exist).To(BeFalse())
			})
		})
		When("the directory does not exist", func() {
			It("should return true", func() {
				exist, err := DirNotExist("path/to/nonexisting/dir")
				Expect(err).NotTo(HaveOccurred())
				Expect(exist).To(BeTrue())
			})
		})
	})
})

var _ = Describe("CheckExistanceOfFile", func() {
	var (
		existingTempFile *os.File
		etfp             string
	)

	BeforeEach(func() {
		var err error

		cwd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		existingTempFile, err = os.CreateTemp(cwd, "test-file-present")
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			err = existingTempFile.Close()
			Expect(err).NotTo(HaveOccurred())
		}()

		etfp = existingTempFile.Name()
	})

	AfterEach(func() {
		err := os.Remove(etfp)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("FilePresent", func() {
		It("should return true if the file exists", func() {
			Expect(FilePresent(etfp)).To(BeTrue())
		})
	})

	Context("FileNotPresent", func() {
		It("should return true if the file does not exist", func() {
			nonExistingFilePath := "path/to/nonexisting/file"
			Expect(FileNotPresent(nonExistingFilePath)).To(BeTrue())
		})
	})
})

var _ = Describe("GetContentLength", func() {
	var (
		ts *httptest.Server
	)

	BeforeEach(func() {
		ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "100")
			w.WriteHeader(http.StatusOK)
		}))
	})

	AfterEach(func() {
		ts.Close()
	})

	It("returns the expected result", func() {
		expectedResult := int64(100)
		result, err := GetContentLength(ts.URL)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(expectedResult))
	})
})

var _ = Describe("GetPublicIp", func() {
	It("should return the public IP", func(ctx SpecContext) {
		ip, err := GetPublicIp()
		Expect(err).To(BeNil())
		Expect(ip).ToNot(BeEmpty())

		ipPattern := `\b(?:\d{1,3}\.){3}\d{1,3}\b`
		Expect(ip).To(MatchRegexp(ipPattern))
	})
})

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
		}, SpecTimeout(time.Second*10))
	})
})
