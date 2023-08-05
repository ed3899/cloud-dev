package utils

import "fmt"

type WithMaskF func(mask int) (GetIp MaskIpL)
type MaskIpL func(ip string) (maskedIp string)

func WithMask(mask int) (MaskIp MaskIpL) {

	MaskIp = func(ip string) (maskedIp string) {
		return fmt.Sprintf("%s/%d", ip, mask)
	}

	return
}
