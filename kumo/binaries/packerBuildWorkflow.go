package binaries

import "github.com/pkg/errors"

func PackerBuildWorkflow(packer *Packer) (err error) {
	var (
		cloud Cloud
	)

	if cloud, err = GetCloud(); err != nil {
		err = errors.Wrap(err, "Error occurred while getting cloud")
		return
	}

	packer.Init(cloud)
}