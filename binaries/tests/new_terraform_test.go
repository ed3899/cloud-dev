package tests

import (
	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewTerraform", func() {
	Context("with a valid manager", func() {
		var (
			exePath  = "/usr/local/bin/terraform"
			_manager = &manager.Manager{
				Tool: iota.Terraform,
				Path: &manager.Path{
					Executable: exePath,
				},
			}
		)

		It("should return a new Terraform", func() {
			terraform, err := binaries.NewTerraform(_manager)
			Expect(err).To(BeNil())
			Expect(terraform).ToNot(BeNil())
			Expect(terraform.Path).To(Equal(exePath))
		})
	})

	Context("with an invalid manager", func() {
		var (
			_manager = &manager.Manager{
				Tool: iota.Packer,
			}
		)

		It("should return an error", func() {
			terraform, err := binaries.NewTerraform(_manager)
			Expect(err).ToNot(BeNil())
			Expect(terraform).To(BeNil())
		})
	})
})