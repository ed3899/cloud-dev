package binz

import "os"

func KumoConfigPresent() bool {
	// Check if the config file exists
	_, err := os.Stat("kumo.config.yaml")

	// Return true if the file exists, otherwise false
	return !os.IsNotExist(err)
}

func KumoConfigNotPresent() bool {
	return !KumoConfigPresent()
}