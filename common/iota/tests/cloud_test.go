package tests

import (
	"github.com/ed3899/kumo/common/iota"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cloud", func() {
	Context("aws", func() {
		It("should return the correct iota", func() {
			Expect(iota.Aws.Iota()).To(Equal(iota.Aws))
		})

		It("should return the correct name", func() {
			Expect(iota.Aws.Name()).To(Equal("aws"))
		})

		It("should return the correct template files", func() {
			Expect(iota.Aws.TemplateFiles()).To(Equal(&iota.TemplateFiles{
				Cloud: "aws.tmpl",
				Base:  "base.tmpl",
			}))
		})
	})
})
