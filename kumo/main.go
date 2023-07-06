package main

import (
	"github.com/ed3899/kumo/utils"
)

func init() {
	host := utils.GetHostSpecs()
	validHost := utils.ValidateHostCompatibility(host)
	packerDependency := utils.DraftPackerDependency(validHost)
	pulumiDependency := utils.DraftPulumiDependency(validHost)
	dependencies := []*utils.Dependency{packerDependency, pulumiDependency}
	utils.DownloadDependencies(dependencies)
}

func main() {

}
