package binz

import (
	"github.com/ed3899/kumo/utils"
)

func GetPulumi(bins *utils.Binaries) *utils.Pulumi {
	return &utils.Pulumi{
		ExecutablePath: bins.Pulumi.Dependency.ExtractionPath,
	}
}
