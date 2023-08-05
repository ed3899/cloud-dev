package cloud

import (
	"path/filepath"

	constants "github.com/ed3899/kumo/0_constants"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Cloud struct {
	Name                         string
	Credentials                  CredentialsI
	AssociatedPackerManifestPath string
}

type PickCloud = func(cloudFromConfig string) (cloud Cloud, err error)

func NewCloud(kumoExecAbsPath string) PickCloud {
	var (
		oopsBuilder = oops.
			Code("new_cloud_failed").
			With("kumoExecAbsPath", kumoExecAbsPath)
	)

	PickCloud := func(cloudFromConfig string) (cloud Cloud, err error) {
		switch cloudFromConfig {
		case constants.AWS:
			cloud = Cloud{
				Name: constants.AWS,
				Credentials: AwsCredentials{
					AccessKeyId:     viper.GetString("AWS.AccessKeyId"),
					SecretAccessKey: viper.GetString("AWS.SecretAccessKey"),
				},
				AssociatedPackerManifestPath: filepath.Join(kumoExecAbsPath, constants.PACKER, constants.AWS, constants.PACKER_MANIFEST),
			}

		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
		}

		return
	}

	return PickCloud
}
