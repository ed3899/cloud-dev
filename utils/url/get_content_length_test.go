package url

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetContentLength", func() {
	Context("when the url is valid", Label("unit"), func() {
		It("should return content length and no error for successful request", func() {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Length", "12345")
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			contentLength, err := GetContentLength(server.URL)
			Expect(err).NotTo(HaveOccurred())
			Expect(contentLength).To(Equal(int64(12345)))
		})
	})

	Context("when the url is not valid", Label("unit"), func() {
		It("should return error for unsuccessful request", func() {
			_, err := GetContentLength("non-existent-url")
			Expect(err).To(HaveOccurred())
		})
	})
})
