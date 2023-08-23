package ip

import (
	"os"
	"regexp"

	"github.com/samber/oops"
)

// Reads the first IP address from a file.
//
// Example:
//	("path/to/file_with_an_ip_address.txt") -> "192.168.1.1"
//
// Assuming the file contains a valid IP address.
func ReadIpFromFile(
	path string,
) (string, error) {
	oopsBuilder := oops.
		Code("ReadIpFromFile").
		In("utils").
		In("ip").
		With("path", path)
	// Define the regular expression pattern for matching an IP address
	ipPattern := "\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b"
	// Compile the regular expression
	ipRegex := regexp.MustCompile(ipPattern)

	// Read the contents of the file
	content, err := os.ReadFile(path)
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "error reading file")
		return "", err
	}

	// Find the first match in the content
	ip := ipRegex.FindString(string(content))
	if len(ip) == 0 {
		err := oopsBuilder.
			Errorf("no valid IP address found in file")
		return "", err
	}

	return ip, nil
}
