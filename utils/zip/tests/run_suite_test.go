package tests

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestZip(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zip Suite", Label("utils", "zip"))
}
