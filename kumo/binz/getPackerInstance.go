package binz

import (
	"github.com/ed3899/kumo/utils"
)

func GetPackerInstance(bins *utils.Binaries) *Packer {
	return &Packer{
		ExecutablePath: bins.Packer.Dependency.ExtractionPath,
	}
}
