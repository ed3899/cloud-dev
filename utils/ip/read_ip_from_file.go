package ip

import (
	"fmt"
	"os"
	"regexp"

	"github.com/samber/oops"
)

func ReadIpFromFile(absPath string) (ip string, err error) {
	var (
		oopsBuilder = oops.
				Code("read_ip_from_file_failed").
				With("absPath", absPath)
		// Define the regular expression pattern for matching an IP address
		ipPattern = "\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b"
		// Compile the regular expression
		ipRegex = regexp.MustCompile(ipPattern)

		content []byte
	)

	// Read the contents of the file
	if content, err = os.ReadFile(absPath); err != nil {
		return "", fmt.Errorf("error reading file: %s", err)
	}

	// Find the first match in the content
	ip = ipRegex.FindString(string(content))

	if len(ip) == 0 {
		err = oopsBuilder.
			Errorf("no valid IP address found in file")
		return
	}

	return
}
