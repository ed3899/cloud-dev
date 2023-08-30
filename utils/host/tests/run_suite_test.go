package tests

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHost(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Host Suite", Label("utils", "host"))
}
