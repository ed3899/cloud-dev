package binz

import (
	"github.com/ed3899/kumo/utils"
)

func GetBinaryInstances(bins *utils.Binaries) (packer *Packer, pulumi *Pulumi) {
	return GetPackerInstance(bins), GetPulumiInstance(bins)
}
