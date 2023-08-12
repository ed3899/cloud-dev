package iota

import "github.com/samber/oops"

func RawCloudToCloudIota(rawCloud string) Cloud {
	var cloud Cloud

	oops.
		In("iota").
		Code("RawCloudToCloudIota").
		Recoverf(
			func() {
				switch rawCloud {
				case "aws":
					cloud = Aws
				default:
					panic(rawCloud)
				}
			},
			"Unknown cloud",
		)

	return cloud
}
