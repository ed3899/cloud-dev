package ip

import "fmt"

// Formats the IP address with the given mask.
//	Example:
// 	("192.168.1.1", 24) -> "192.168.1.1/24"
func MaskIp(
	ip string,
	mask int,
) string {
	return fmt.Sprintf("%s/%d", ip, mask)
}
