package binz

import (
	"log"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

func GetBinaryInstances(bins *utils.Binaries) (packer *Packer, pulumi *Pulumi) {
	packer, err := GetPackerInstance(bins)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Packer instance")
		log.Fatal(err)
	}

	pulumi, err = GetPulumiInstance(bins)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Pulumi instance")
		log.Fatal(err)
	}

	return packer, pulumi
}
