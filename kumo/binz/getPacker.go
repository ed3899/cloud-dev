package binz

import (
	"github.com/ed3899/kumo/utils"
)

func GetPacker(bins *utils.Binaries) *utils.Packer {
	return &utils.Packer{
		ExecutablePath: bins.Packer.Dependency.ExtractionPath,
	}
}