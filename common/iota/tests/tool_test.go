package tests

import (
	"github.com/ed3899/kumo/common/iota"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tool", func() {
	Context("terraform", func() {
		It("should return the correct name", func() {
			Expect(iota.Terraform.Name()).To(Equal("terraform"))
		})

		It("should return the correct vars name", func() {
			Expect(iota.Terraform.VarsName()).To(Equal(".auto.tfvars"))
		})

		It("should return the correct version", func() {
			Expect(iota.Terraform.Version()).To(Equal("1.5.5"))
		})

		It("should return the correct pluging path environment variable", func() {
			Expect(iota.Terraform.PluginPathEnvironmentVariable()).To(Equal("TF_PLUGIN_CACHE_DIR"))
		})

		It("should return the correct plugin dir", func() {
			Expect(iota.Terraform.PluginDir()).To(Equal(".terraform"))
		})
	})

	Context("packer", func() {
		It("should return the correct name", func() {
			Expect(iota.Packer.Name()).To(Equal("packer"))
		})

		It("should return the correct vars name", func() {
			Expect(iota.Packer.VarsName()).To(Equal(".auto.pkrvars.hcl"))
		})

		It("should return the correct version", func() {
			Expect(iota.Packer.Version()).To(Equal("1.9.2"))
		})

		It("should return the correct pluging path environment variable", func() {
			Expect(iota.Packer.PluginPathEnvironmentVariable()).To(Equal("PACKER_PLUGIN_PATH"))
		})

		It("should return the correct plugin dir", func() {
			Expect(iota.Packer.PluginDir()).To(Equal("plugins"))
		})
	})
})
