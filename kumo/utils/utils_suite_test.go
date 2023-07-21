package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ed3899/kumo/host"
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

var _ = Describe("GetDependenciesDirName", func() {
	It("should return the correct dependencies directory name", func() {
		Expect(GetDependenciesDirName()).To(Equal("deps"))
	})
})

var _ = Describe("GetDependencyURL", func() {
	isValidURL := func(input string) bool {
		_, err := url.ParseRequestURI(input)
		return err == nil
	}

	type TestCase struct {
		specs       *host.Specs
		name        string
		expectedURL string
	}

	latestPackerVersion := GetLatestPackerVersion()
	latestTerraformVersion := GetLatestTerraformVersion()

	DescribeTable("returns the correct dependency URL with valid format", func(tc *TestCase) {
		url, err := GetDependencyURL(tc.name, tc.specs)
		Expect(err).ToNot(HaveOccurred())
		Expect(url).To(Equal(tc.expectedURL))
		Expect(isValidURL(url)).To(BeTrue())
	},
		Entry("when crafting for packer on windows amd64", &TestCase{
			specs: &host.Specs{
				OS:   "windows",
				ARCH: "amd64",
			},
			name:        "packer",
			expectedURL: fmt.Sprintf("https://releases.hashicorp.com/packer/%s/packer_%s_windows_amd64.zip", latestPackerVersion, latestPackerVersion),
		}),
		Entry("when crafting for packer on windows 386", &TestCase{
			specs: &host.Specs{
				OS:   "windows",
				ARCH: "386",
			},
			name:        "packer",
			expectedURL: fmt.Sprintf("https://releases.hashicorp.com/packer/%s/packer_%s_windows_386.zip", latestPackerVersion, latestPackerVersion),
		}),
		Entry("when crafting for packer on macos amd64", &TestCase{
			specs: &host.Specs{
				OS:   "darwin",
				ARCH: "amd64",
			},
			name:        "packer",
			expectedURL: fmt.Sprintf("https://releases.hashicorp.com/packer/%s/packer_%s_darwin_amd64.zip", latestPackerVersion, latestPackerVersion),
		}),
		Entry("when crafting for packer on macos arm64", &TestCase{
			specs: &host.Specs{
				OS:   "darwin",
				ARCH: "arm64",
			},
			name:        "packer",
			expectedURL: fmt.Sprintf("https://releases.hashicorp.com/packer/%s/packer_%s_darwin_arm64.zip", latestPackerVersion, latestPackerVersion),
		}),
		Entry("when crafting for terraform on macos arm64", &TestCase{
			specs: &host.Specs{
				OS:   "darwin",
				ARCH: "arm64",
			},
			name:        "terraform",
			expectedURL: fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_darwin_arm64.zip", latestTerraformVersion, latestTerraformVersion),
		}),
		Entry("when crafting for terraform on macos amd64", &TestCase{
			specs: &host.Specs{
				OS:   "darwin",
				ARCH: "amd64",
			},
			name:        "terraform",
			expectedURL: fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_darwin_amd64.zip", latestTerraformVersion, latestTerraformVersion),
		}),
	)
})

var _ = Describe("GetLastBuiltAmiId", func() {
	var tmpFile *os.File
	var expectedAmiId = "67890"

	BeforeEach(func() {
		var err error
		tmpFile, err = os.CreateTemp("", "test_get_last_built_ami_id")
		Expect(err).ToNot(HaveOccurred())

		jsonData := fmt.Sprintf(`{
			"builds": [
							{
											"packer_run_uuid": "abc123",
											"artifact_id": "ami:12345"
							},
							{
											"packer_run_uuid": "def456",
											"artifact_id": "ami:%s"
							}
			],
			"last_run_uuid": "def456"
		}`, expectedAmiId)

		_, err = tmpFile.WriteString(jsonData)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	})

	It("returns the last built AMI ID", func() {
		// Run the test using the temporary packer manifest file
		amiId, err := GetLastBuiltAmiId(tmpFile.Name())
		Expect(err).ToNot(HaveOccurred())

		Expect(amiId).To(Equal(expectedAmiId))
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
