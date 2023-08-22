package url

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Download", func() {
	var (
		server *httptest.Server
		url    string
		path   string
	)

	BeforeEach(func() {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "9")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test data"))
		}))

		url = server.URL
		path = filepath.Join(os.TempDir(), "test_download_file.txt")
	})

	AfterEach(func() {
		server.Close()
		_ = os.Remove(path)
	})

	It("should download content from URL and save it to the file", func() {
		var bytesDownloaded int
		bytesDownloadedChan := make(chan int, 1024)

		err := Download(url, path, bytesDownloadedChan)
		close(bytesDownloadedChan)
		for b := range bytesDownloadedChan {
			bytesDownloaded += b
		}

		Expect(err).NotTo(HaveOccurred())
		Expect(bytesDownloaded).To(Equal(9)) // Length of "test data"
		Expect(path).To(BeAnExistingFile())

		content, err := os.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(content)).To(Equal("test data"))
	})
})
