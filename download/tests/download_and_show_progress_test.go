package tests

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/ed3899/kumo/download"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vbauerster/mpb/v8"
)

var _ = Describe("DownloadAndShowProgress", func() {
	Context("with a valid url", func() {
		var (
			mockDownload *download.Download

			cwd string
			err error
		)

		BeforeEach(func() {
			cwd, err = os.Getwd()
			Expect(err).To(BeNil())

			// The bar won't display any progress because content length was not added. We only care about the download
			mockDownload = &download.Download{
				Url: "https://releases.hashicorp.com/packer/1.9.4/packer_1.9.4_darwin_amd64.zip",
				Path: &download.Path{
					Zip: filepath.Join(cwd, "packer.zip"),
				},
				Progress: mpb.New(mpb.WithWaitGroup(&sync.WaitGroup{}), mpb.WithAutoRefresh(), mpb.WithWidth(64)),
				Bar:      &download.Bar{},
			}
		})

		AfterEach(func() {
			Expect(os.Remove(mockDownload.Path.Zip)).To(Succeed())
		})

		It("should download and show progress", Label("integration"), func() {
			Expect(mockDownload.DownloadAndShowProgress()).To(Succeed())
		})
	})

	Context("with an invalid url", func() {
		var (
			mockDownload *download.Download

			cwd string
			err error
		)

		BeforeEach(func() {
			cwd, err = os.Getwd()
			Expect(err).To(BeNil())

			// The bar won't display any progress because content length was not added. We only care about the download
			mockDownload = &download.Download{
				Url: "invalid_url",
				Path: &download.Path{
					Zip: filepath.Join(cwd, "packer.zip"),
				},
				Progress: mpb.New(mpb.WithWaitGroup(&sync.WaitGroup{}), mpb.WithAutoRefresh(), mpb.WithWidth(64)),
				Bar:      &download.Bar{},
			}
		})

		AfterEach(func() {
			Expect(os.RemoveAll(mockDownload.Path.Zip)).To(Succeed())
		})

		It("should handle download error", Label("integration"), func() {
			Expect(mockDownload.DownloadAndShowProgress()).ToNot(Succeed())
		})
	})
})
