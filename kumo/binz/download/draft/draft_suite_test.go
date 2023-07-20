package draft_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDraft(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Draft Suite")

}
