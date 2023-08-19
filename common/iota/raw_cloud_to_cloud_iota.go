package iota

import (
	"log"

	"github.com/samber/oops"
)

func CloudIota(
	rawCloudFromConfig string,
) Cloud {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Code("CloudIota").
		With("rawCloud", rawCloudFromConfig)

	switch rawCloudFromConfig {
	case "aws":
		return Aws

	default:
		err := oopsBuilder.
			Errorf("unknown cloud: %#v", rawCloudFromConfig)

		log.Fatalf("%+v", err)

		return -1
	}
}
