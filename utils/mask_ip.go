package utils

import "fmt"

// Return the given IP address with the given mask
//
// Example:
//
//	("0.0.0.0", 32) -> "0.0.0.0/32"
func MaskIp(ip string, mask int) (maskedIp string) {
	return fmt.Sprintf("%s/%d", ip, mask)
}