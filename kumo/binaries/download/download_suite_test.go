package download

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ed3899/kumo/binz/download/draft"
	"github.com/ed3899/kumo/binz/download/mocks"
	"github.com/ed3899/kumo/binz/download/progressBar"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var (
	mockServer   *httptest.Server
	tempDir      string
	zipPath      string
	downloads    chan *progressBar.DownloadResult
	responseBody string
	downloadBar  draft.ProgressBar
	ctrl         *gomock.Controller
)

func TestDraft(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Draft Suite")
}

var _ = BeforeSuite(func() {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}))
})

var _ = AfterSuite(func() {
	mockServer.Close()
})

var _ = BeforeEach(func() {
	ctrl = gomock.NewController(GinkgoT())
	tempDir, _ = os.MkdirTemp("", "test")
	zipPath = filepath.Join(tempDir, "test.zip")
	downloads = make(chan *progressBar.DownloadResult, 1)
	downloadBar = mocks.NewMockProgressBar(ctrl)

	downloadBar.(*mocks.MockProgressBar).EXPECT().IncrBy(gomock.Any()).AnyTimes()
})

var _ = AfterEach(func() {
	os.RemoveAll(tempDir)
})

var _ = Describe("Download single dependency", func() {

	Context("when download is a valid url", func() {
		responseBody = "Hello, World!"

		It("should download and save the file successfully", func() {
			dep := &draft.Dependency{
				URL:          mockServer.URL,
				DownloadPath: zipPath,
				DownloadBar:  downloadBar,
			}

			DownloadDependency(dep, downloads)

			result := <-downloads
			Expect(result.Err).To(BeNil())

			content, _ := os.ReadFile(zipPath)
			Expect(string(content)).To(Equal(responseBody))
		})
	})

	Context("when download is with an invalid url", func() {
		It("should return an error", func() {
			dep := &draft.Dependency{
				URL:          "wrong-url",
				DownloadPath: zipPath,
				DownloadBar:  downloadBar,
			}

			DownloadDependency(dep, downloads)

			result := <-downloads
			Expect(result.Err).To(HaveOccurred())
		})
	})
})
