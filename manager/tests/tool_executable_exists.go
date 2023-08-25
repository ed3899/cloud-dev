package tests

import (
	"os"

	"github.com/ed3899/kumo/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ToolExecutableExists", func() {
	Context("with an executable path that exists", func() {
		var (
			tempFile *os.File
			_manager *manager.Manager
		)

		BeforeEach(func() {
			cwd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			tempFile, err = os.CreateTemp(cwd, "tempfile")
			defer tempFile.Close()
			Expect(err).ToNot(HaveOccurred())

			_manager = &manager.Manager{
				Path: &manager.Path{
					Executable: tempFile.Name(),
				},
			}
		})

		AfterEach(func() {
			Expect(os.Remove(tempFile.Name())).To(Succeed())
		})

		It("should return true", Label("unit"), func() {
			Expect(_manager.ToolExecutableExists()).To(BeTrue())
		})
	})

	Context("with an executable path that does not exist", func() {
		var (
			_manager *manager.Manager
		)

		BeforeEach(func() {
			_manager = &manager.Manager{
				Path: &manager.Path{
					Executable: "does-not-exist",
				},
			}
		})

		It("should return false", Label("unit"), func() {
			Expect(_manager.ToolExecutableExists()).To(BeFalse())
		})
	})
})
