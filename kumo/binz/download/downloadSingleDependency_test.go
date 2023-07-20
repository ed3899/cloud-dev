package download

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ed3899/kumo/binz/download/draft"
	"github.com/ed3899/kumo/binz/download/progressBar"
	"github.com/stretchr/testify/assert"
)

func TestDownload_Success(t *testing.T) {
	// Create a mock server to simulate the HTTP response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer mockServer.Close()

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a path for the zip file
	zipPath := filepath.Join(tempDir, "test.zip")

	// Create a Dependency instance
	dep := &draft.Dependency{
		URL:     mockServer.URL,
		ZipPath: zipPath,
	}

	// Create a channel to capture the download result
	downloads := make(chan *progressBar.DownloadResult, 1)

	// Execute the Download function
	Download(dep, downloads)

	// Expect the channel to receive the download result
	result := <-downloads
	assert.NoError(t, result.Err)

	// Verify that the zip file was downloaded successfully
	data, err := os.ReadFile(zipPath)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(data))
}

// func TestDownload_Failure(t *testing.T) {
// 	// Create a Dependency instance with an invalid URL
// 	dep := &draft.Dependency{
// 		URL:     "https://example.com/nonexistent-file",
// 		ZipPath: "/path/to/nonexistent-directory/test.zip",
// 	}

// 	// Create a channel to capture the download result
// 	downloads := make(chan *progressBar.DownloadResult, 1)

// 	// Execute the Download function
// 	Download(dep, downloads)

// 	// Expect the channel to receive an error in the download result
// 	result := <-downloads
// 	assert.Error(t, result.Err)
// }

// func TestMain(m *testing.M) {
// 	// Setup any necessary test fixtures or configurations
// 	// ...

// 	// Run tests
// 	code := m.Run()

// 	// Teardown any resources if required
// 	// ...

// 	os.Exit(code)
// }
