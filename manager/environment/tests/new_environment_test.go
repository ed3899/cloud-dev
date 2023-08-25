package tests

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager/environment"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewEnvironment", func() {
	Context("when the tool is unknown", func() {
		It("returns an error", func() {
			_, err := environment.NewEnvironment(
				-1,
				iota.Aws,
				"",
			)
			Expect(err).To(HaveOccurred())
		})
	})
})
