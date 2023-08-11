package cloud

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func NewCloud(
	options ...Option,
) (
	cloud *Cloud,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("NewCloud").
				With("opts", options)

		opt Option
	)

	cloud = &Cloud{}
	for _, opt = range options {
		if err = opt(cloud); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Option %v", opt)
			return
		}
	}

	return
}

func WithKind(
	cloudFromConfig string,
) (
	option Option,
) {
	var (
		oopsBuilder = oops.
			Code("WithKind").
			With("cloudFromConfig", cloudFromConfig)
	)

	option = func(cloud *Cloud) (err error) {
		switch cloudFromConfig {
		case constants.AWS:
			cloud.kind = constants.Aws
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
			return
		}

		return
	}

	return
}

func WithName(
	cloudFromConfig string,
) (
	option Option,
) {
	var (
		oopsBuilder = oops.
			Code("WithName").
			With("cloudFromConfig", cloudFromConfig)
	)

	option = func(cloud *Cloud) (err error) {
		switch cloudFromConfig {
		case constants.AWS:
			cloud.name = constants.AWS
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
			return
		}

		return
	}

	return
}

func WithCredentials(
	cloudFromConfig string,
) (
	option Option,
) {
	var (
		oopsBuilder = oops.
			Code("WithCredentials").
			With("cloudFromConfig", cloudFromConfig)
	)

	option = func(cloud *Cloud) (err error) {
		switch cloudFromConfig {
		case constants.AWS:
			cloud.Credentials = AwsCredentials{
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

	return
}

func WithPackerManifestPath(
	cloudFromConfig,
	kumoExecAbsPath string,
) (
	option Option,
) {
	var (
		oopsBuilder = oops.
				Code("WithPackerManifestPath").
				With("cloudFromConfig", cloudFromConfig)

		packerManifestPath = func(cloudName string) (path PackerManifestPath) {
			path = PackerManifestPath(
				filepath.Join(
					kumoExecAbsPath,
					constants.PACKER,
					cloudName,
					constants.PACKER_MANIFEST,
				),
			)

			return
		}
	)

	option = func(cloud *Cloud) (err error) {
		switch cloudFromConfig {
		case constants.AWS:
			cloud.packerManifestPath = packerManifestPath(constants.AWS)
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
			return
		}

		return
	}

	return
}

func (c *Cloud) Kind() (cloudKind constants.CloudKind) {
	cloudKind = c.kind

	return
}

func (c *Cloud) Name() (cloudName CloudName) {
	cloudName = c.name

	return
}

func (c *Cloud) PackerManifestPath() (packerManifestPath PackerManifestPath) {
	packerManifestPath = c.packerManifestPath

	return
}

type Option func(*Cloud) error

type CloudName string

func (c CloudName) String() (s string) {
	s = string(c)

	return
}

type PackerManifestPath string

func (p PackerManifestPath) String() (s string) {
	s = string(p)

	return
}

type Cloud struct {
	kind               constants.CloudKind
	name               CloudName
	packerManifestPath PackerManifestPath
	Credentials        interfaces.Credentials
}

type PackerManifestPathGetter interface {
	PackerManifestPath() PackerManifestPath
}

type CloudI interface {
	interfaces.KindGetter[constants.CloudKind]
	interfaces.NameGetter[CloudName]
	PackerManifestPathGetter
}
