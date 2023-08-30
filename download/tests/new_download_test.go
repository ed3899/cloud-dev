package tests

import (
	"fmt"
	"net/url"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ed3899/kumo/download"
	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager"
)

var _ = Describe("NewDownload", func() {
	var (
		mockManager      *manager.Manager
		zipPathSubstring string
		exePathSubstring string
		urlSubstring     string
		err              error
	)

	BeforeEach(func() {
		mockManager = &manager.Manager{
			Tool: iota.Packer,
		}

		zipPathSubstring = filepath.Join(
			iota.Dependencies.Name(),
			fmt.Sprintf("%s.zip", iota.Packer.Name()),
		)

		exePathSubstring = filepath.Join(
			iota.Dependencies.Name(),
			iota.Packer.Name(),
			fmt.Sprintf("%s.exe", iota.Packer.Name()),
		)

		urlSubstring, err = url.JoinPath(
			iota.Packer.Name(),
			iota.Packer.Version(),
		)
		Expect(err).To(BeNil())
	})

	It("should create a new download struct", Label("integration"), func() {
		_download, err := download.NewDownload(mockManager)
		Expect(err).To(BeNil())
		Expect(_download.Name).To(Equal(mockManager.Tool.Name()))
		Expect(_download.Path.Zip).To(ContainSubstring(zipPathSubstring))
		Expect(_download.Path.Executable).To(ContainSubstring(exePathSubstring))
		Expect(_download.Url).To(ContainSubstring(urlSubstring))
		Expect(_download.ContentLength).To(BeNumerically(">", 0))
		Expect(_download.Progress).ToNot(BeNil())
		Expect(_download.Bar).ToNot(BeNil())
	})
})
