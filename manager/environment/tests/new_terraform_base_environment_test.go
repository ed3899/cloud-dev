package tests

import (
	"github.com/ed3899/kumo/manager/environment"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewTerraformBaseEnvironment", func() {
	It("should return a valid TerraformBaseEnvironment", func() {
		_environment := environment.NewTerraformBaseEnvironment()

		Expect(_environment).ToNot(BeNil())
		Expect(_environment.Required.ALLOWED_IP).ToNot(BeEmpty())
	})
})
