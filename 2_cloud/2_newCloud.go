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

func NewCloud(cloudFromConfig string, kumoExecAbsPath string) (cloud Cloud, err error) {
	var (
		oopsBuilder = oops.
			Code("new_cloud_failed").
			With("cloudFromConfig", cloudFromConfig)
	)

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
