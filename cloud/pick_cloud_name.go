package cloud

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/mo"
	"github.com/samber/oops"
)

func PickCloudName(
	cloudIota iota.Cloud,
) mo.Result[string] {
	oopsBuilder := oops.
		Code("CloudNameWith").
		With("cloudIota", cloudIota)

	switch cloudIota {
	case iota.Aws:
		return mo.Ok("aws")

	default:
		err := oopsBuilder.Errorf(
			"Unknown cloud '%#v'",
			cloudIota,
		)
		return mo.Err[string](err)
	}
}
