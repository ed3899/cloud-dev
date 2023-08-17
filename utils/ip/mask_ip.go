package ip

import "fmt"

func MaskIp(
	ip string,
	mask int,
) string {
	return fmt.Sprintf("%s/%d", ip, mask)
}
