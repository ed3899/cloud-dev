package ip

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MaskIp", func() {
	It("should format the IP address with the given mask", Label("unit"), func() {
		ip := "192.168.1.1"
		mask := 24
		expectedResult := fmt.Sprintf("%s/%d", ip, mask)

		result := MaskIp(ip, mask)

		Expect(result).To(Equal(expectedResult))
	})
})
