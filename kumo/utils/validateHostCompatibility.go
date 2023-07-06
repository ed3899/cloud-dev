package utils

import "log"

func ValidateHostCompatibility(s Specs) Specs {
	switch s.OS {
	case "windows":
		switch s.ARCH {
		case "386":
		case "amd64":
		default:
			log.Fatalf("Looks like your operative systems architecture is not supported :/")
		}
	default:
		log.Fatalf("Looks like your operative system is not supported :/")
	}

	return s
}