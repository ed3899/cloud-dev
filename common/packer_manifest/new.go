package packer_manifest

import (
	cloud_interfaces "github.com/ed3899/kumo/common/cloud/interfaces"
	tool_constants "github.com/ed3899/kumo/common/tool/constants"
	"github.com/samber/oops"
)

type Manifest struct {
	lastBuiltAmiId string
}

func New(cloud cloud_interfaces.Cloud) (manifest *Manifest, err error) {
	var (
		oopsBuilder = oops.
			Code("new_manifest_failed")

		absPath string
	)


}
