package download

import (
	"fmt"
	"net/url"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

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
		download, err := NewDownload(mockManager)
		Expect(err).To(BeNil())
		Expect(download.Name).To(Equal(mockManager.Tool.Name()))
		Expect(download.Path.Zip).To(ContainSubstring(zipPathSubstring))
		Expect(download.Path.Executable).To(ContainSubstring(exePathSubstring))
		Expect(download.Url).To(ContainSubstring(urlSubstring))
		Expect(download.ContentLength).To(BeNumerically(">", 0))
		Expect(download.Progress).ToNot(BeNil())
		Expect(download.Bar).ToNot(BeNil())
	})
})
