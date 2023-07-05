package main

import (
	"fmt"
	"os"
	"os/exec"
)

// TODO how to embed binaries
// TODO how to access folders
// TODO how to set env vars
// TODO Detect and install packer based on os and arc

func main() {
	// Set packer plugin path
	env := os.Environ()
	env = append(env, "PACKER_PLUGIN_PATH=/path/to/plugin/dir")

	cmd := exec.Command("./packer/windows/packer_1.9.1_windows_amd64/packer.exe", "--version")
	cmd.Env = env

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("there was an error while calling packer:", err)
		return
	}

	fmt.Println(string(output))
}