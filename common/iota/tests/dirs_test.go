package tests

import (
	"github.com/ed3899/kumo/common/iota"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dirs", func() {
	Context("dependencies", func() {
		It("should return the correct name", func() {
			Expect(iota.Dependencies.Name()).To(Equal("dependencies"))
		})
	})

	Context("templates", func() {
		It("should return the correct name", func() {
			Expect(iota.Templates.Name()).To(Equal("templates"))
		})
	})
})
