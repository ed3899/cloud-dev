package cloud

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

type CloudManager struct {
	cloud iota.Cloud
}

func (c CloudManager) PickCloud(rawCloud string) CloudManager {
	err := oops.
		In("cloud").
		Tags("CloudManager").
		Code("PickCloud").Recoverf(
		func() {},
	)

	switch rawCloud {
	case "aws":
		c.cloud = iota.Aws
	default:
		panic("Unknown cloud")
	}

	return c
}

func (c CloudManager) Clone() CloudManager {
	return CloudManager{
		cloud: c.cloud,
	}
}
