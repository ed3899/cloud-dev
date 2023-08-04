package cloud

import (
	"path/filepath"

	"github.com/samber/oops"
	"github.com/spf13/viper"
)

const (
	PACKER          = "packer"
	PACKER_MANIFEST = "manifest.json"
	AWS             = "aws"
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
	case AWS:
		cloud = Cloud{
			Name: AWS,
			Credentials: AwsCredentials{
				AccessKeyId:     viper.GetString("AWS.AccessKeyId"),
				SecretAccessKey: viper.GetString("AWS.SecretAccessKey"),
			},
			AssociatedPackerManifestPath: filepath.Join(kumoExecAbsPath, PACKER, AWS, PACKER_MANIFEST),
		}

	default:
		err = oopsBuilder.
			Wrapf(err, "Cloud %s is not supported", cloudFromConfig)

	}

	return
}
