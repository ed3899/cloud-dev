package binz

import (
	"github.com/ed3899/kumo/utils"
)

func GetPulumiInstance(bins *utils.Binaries) *Pulumi {
	return &Pulumi{
		ExecutablePath: bins.Pulumi.Dependency.ExtractionPath,
	}
}
