package main

import (
	// "fmt"
	// "math/rand"
	// "sync"
	// "time"

	"github.com/ed3899/kumo/utils"

	// "github.com/vbauerster/mpb/v8"
	// "github.com/vbauerster/mpb/v8/decor"
)

// TODO how to embed binaries
// TODO how to access folders
// TODO how to set env vars
// TODO Detect and install packer based on os and arc

func init() {
	hs := utils.GetHostSpecs()
	vh := utils.ValidateHostCompatibility(hs)
	packerUrl := utils.GetPackerUrl(vh)
	pulumiUrl := utils.GetPulumiUrl(vh)
	urls := []*utils.ZipExecutableRef{packerUrl, pulumiUrl}
	utils.DownloadPackages(urls)
}

func main() {

}