package tests

import (
	"fmt"

	"github.com/ed3899/kumo/utils/ip"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MaskIp", func() {
	It("should format the IP address with the given mask", Label("unit"), func() {
		_ip := "192.168.1.1"
		mask := 24
		expectedResult := fmt.Sprintf("%s/%d", _ip, mask)

		result := ip.MaskIp(_ip, mask)

		Expect(result).To(Equal(expectedResult))
	})
})
