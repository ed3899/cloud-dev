package utils

// Return the IP to be used based on the public IP provided.
// If the public IP is empty, returns "0.0.0.0"
func PickIpToBeUsed(publicIp string) (ipToBeUsed string) {
	switch {
	case publicIp == "":
		return "0.0.0.0"
	default:
		return publicIp
	}
}
