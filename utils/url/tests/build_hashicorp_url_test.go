package tests

import (
	"fmt"

	"github.com/ed3899/kumo/utils/url"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BuildHashicorpUrl", func() {
	It("should build the correct URL", Label("unit"), func() {
		name := "terraform"
		version := "0.12.0"
		os := "linux"
		arch := "amd64"

		expectedURL := fmt.Sprintf("https://releases.hashicorp.com/%s/%s/%s_%s_%s_%s.zip", name, version, name, version, os, arch)

		_url := url.BuildHashicorpUrl(name, version, os, arch)
		Expect(_url).To(Equal(expectedURL))
	})
})
