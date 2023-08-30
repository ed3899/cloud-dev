package tests

import (
	"github.com/ed3899/kumo/common/iota"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudIota", func() {
	Context("aws", func() {
		It("should return the correct iota", func() {
			Expect(iota.CloudIota("aws")).To(Equal(iota.Aws))
		})
	})
})
