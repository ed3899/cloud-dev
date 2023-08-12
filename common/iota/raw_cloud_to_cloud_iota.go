package iota

import "github.com/samber/oops"

func RawCloudToCloudIota(rawCloud string) (Cloud, error) {
	oopsBuilder := oops.
		In("iota").
		Code("RawCloudToCloudIota").
		With("rawCloud", rawCloud)

	switch rawCloud {
	case "aws":
		return Aws, nil
	default:
		err := oopsBuilder.
			Errorf("unknown cloud: %#v", rawCloud)
		return -1, err
	}
}
