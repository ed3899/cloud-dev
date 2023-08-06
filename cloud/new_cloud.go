package cloud

import (
	"path/filepath"

	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Cloud struct {
	kind               constants.CloudKind
	name               string
	credentials        CredentialsI
	packerManifestPath string
}

type Option func(Cloud) (Cloud, error)

func WithKind(cloudFromConfig string) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithKind").
			With("cloudFromConfig", cloudFromConfig)
	)

	return func(c Cloud) (cloud Cloud, err error) {
		switch cloudFromConfig {
		case constants.AWS:
			c.kind = constants.Aws
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
			return
		}

		return
	}
}

func WithName(cloudFromConfig string) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithName").
			With("cloudFromConfig", cloudFromConfig)
	)

	return func(c Cloud) (cloud Cloud, err error) {
		switch cloudFromConfig {
		case constants.AWS:
			c.name = constants.AWS
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
			return
		}

		return
	}
}

func WithCredentials(cloudFromConfig string) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithCredentials").
			With("cloudFromConfig", cloudFromConfig)
	)

	return func(c Cloud) (cloud Cloud, err error) {
		switch cloudFromConfig {
		case constants.AWS:
			c.credentials = AwsCredentials{
				AccessKeyId:     viper.GetString("AWS.AccessKeyId"),
				SecretAccessKey: viper.GetString("AWS.SecretAccessKey"),
			}
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
			return
		}

		return
	}
}

func WithPackerManifestPath(cloudFromConfig, kumoExecAbsPath string) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithPackerManifestPath").
			With("cloudFromConfig", cloudFromConfig)
	)

	return func(c Cloud) (cloud Cloud, err error) {
		switch cloudFromConfig {
		case constants.AWS:
			c.packerManifestPath = filepath.Join(kumoExecAbsPath, constants.PACKER, constants.AWS, constants.PACKER_MANIFEST)
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
			return
		}

		return
	}
}

func NewCloud(opts ...Option) (cloud Cloud, err error) {
	var (
		oopsBuilder = oops.
			Code("new_cloud_failed").
			With("opts", opts)
	)

	cloud = Cloud{}
	for _, opt := range opts {
		if cloud, err = opt(cloud); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Option %v", opt)
			return
		}
	}

	return
}
