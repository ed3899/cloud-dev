package url

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Download", func() {
	var (
		content = "test data"

		server *httptest.Server
		path   string
	)

	BeforeEach(func() {
		// Create a test server that returns the content
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(content))
		}))

		// Create a temporary file to download to
		cwd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		path = filepath.Join(cwd, "test_download_file.txt")
	})

	AfterEach(func() {
		server.Close()
		_ = os.Remove(path)
	})

	Context("when the URL is valid", func() {
		It("should download content from URL and save it to the file", func() {
			bytesDownloaded := 0
			bytesDownloadedChan := make(chan int, 1024)

			wg := &sync.WaitGroup{}
			wg.Add(1)

			go func() {
				defer wg.Done()

				for b := range bytesDownloadedChan {
					bytesDownloaded += b
				}
			}()

			err := Download(server.URL, path, bytesDownloadedChan)
			close(bytesDownloadedChan)
			wg.Wait()

			Expect(err).NotTo(HaveOccurred())
			Expect(bytesDownloaded).To(Equal(len(content))) // Length of "test data"
			Expect(path).To(BeAnExistingFile())

			fileContent, err := os.ReadFile(path)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(fileContent)).To(Equal(content))
		})
	})

	Context("when the URL is invalid", func() {
		It("should return an error", func() {
			err := Download("invalid-url", path, nil)
			Expect(err).To(HaveOccurred())
		})
	})
})
