package cloud

import (
	"path/filepath"

	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Cloud struct {
	Kind                         constants.CloudKind
	Name                         string
	Credentials                  CredentialsI
	AssociatedPackerManifestPath string
}

func NewCloud(kumoExecAbsPath string) (PickCloud PickCloud) {
	var (
		oopsBuilder = oops.
			Code("new_cloud_failed").
			With("kumoExecAbsPath", kumoExecAbsPath)
	)

	PickCloud = func(cloudFromConfig string) (cloud Cloud, err error) {
		switch cloudFromConfig {
		case constants.AWS:
			cloud = Cloud{
				Kind: constants.Aws,
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

	return
}

type PickCloud = func(cloudFromConfig string) (cloud Cloud, err error)
