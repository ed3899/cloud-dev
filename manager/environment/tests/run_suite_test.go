package tests

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestManagerEnvironment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Manager Environment Suite", Label("manager", "environment"))
}
