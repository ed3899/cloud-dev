package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetContentLength(t *testing.T) {
	expectedResult := int64(100)

	// ts is an httptest.Server which will be used as a mock server for testing.
	// It handles incoming requests with a handler that sets the Content-Length header to 100 and sends a 200 OK response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Length", "100")
	}))
	defer ts.Close()

	// Run GetContentLength() with the mock server.
	result, err := GetContentLength(ts.URL)

	if err != nil {
		t.Errorf("GetContentLength() failed, expected %v, got %v", expectedResult, result)
	}
}
