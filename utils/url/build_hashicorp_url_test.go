package url

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BuildHashicorpUrl", func() {
	It("should build the correct URL", func() {
		name := "terraform"
		version := "0.12.0"
		os := "linux"
		arch := "amd64"

		expectedURL := fmt.Sprintf("https://releases.hashicorp.com/%s/%s/%s_%s_%s_%s.zip", name, version, name, version, os, arch)

		url := BuildHashicorpUrl(name, version, os, arch)
		Expect(url).To(Equal(expectedURL))
	})
})
