package url

import "fmt"

// Returns a Hashicorp URL.
//
// Example:
// 	("packer", "1.7.4", "windows", "amd64") -> "https://releases.hashicorp.com/packer/1.7.4/packer_1.7.4_windows_amd64.zip"
func BuildHashicorpUrl(
	name,
	version,
	os,
	arch string,
) string {
	return fmt.Sprintf("https://releases.hashicorp.com/%s/%s/%s_%s_%s_%s.zip", name, version, name, version, os, arch)
}
