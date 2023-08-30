package tests

import (
	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewPacker", func() {
	Context("with a valid manager", func() {
		var (
			exePath  = "/usr/local/bin/packer"
			_manager = &manager.Manager{
				Tool: iota.Packer,
				Path: &manager.Path{
					Executable: exePath,
				},
			}
		)

		It("should return a new Packer", func() {
			packer, err := binaries.NewPacker(_manager)
			Expect(err).To(BeNil())
			Expect(packer).ToNot(BeNil())
			Expect(packer.Path).To(Equal(exePath))
		})
	})

	Context("with an invalid manager", func() {
		var (
			_manager = &manager.Manager{
				Tool: iota.Terraform,
			}
		)

		It("should return an error", func() {
			packer, err := binaries.NewPacker(_manager)
			Expect(err).ToNot(BeNil())
			Expect(packer).To(BeNil())
		})
	})
})
