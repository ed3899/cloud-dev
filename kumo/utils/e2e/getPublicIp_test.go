package utils

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicIpIntegration(t *testing.T) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP status code 200, got %d", resp.StatusCode)

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	ip := string(bytesResp)

	assert.NotEmpty(t, ip, "Expected non-empty IP address, got %s", ip)
}
