package binz

import (
	"github.com/ed3899/kumo/utils"
)

func GetPulumi(bins *utils.Binaries) *Pulumi {
	return &Pulumi{
		ExecutablePath: bins.Pulumi.Dependency.ExtractionPath,
	}
}
