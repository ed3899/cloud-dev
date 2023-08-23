package packer_manifest_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPackerManifest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PackerManifest Suite", Label("utils", "packer_manifest"))
}
