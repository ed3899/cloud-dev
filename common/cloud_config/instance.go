package cloud_config

import (
	"github.com/ed3899/kumo/common/cloud_config/aws"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Credentials interface {
	Set() error
	Unset() error
}

type Cloud struct {
	name        string
	kind       Kind
	Credentials Credentials
}

func New(cloud string) (cloudConfig *Cloud, err error) {
	var (
		oopsBuilder = oops.
				Code("new_cloud_deployment_failed").
				With("cloud", cloud)
		awsName = AWS_NAME
	)

	switch cloud {
	case awsName:
		cloudConfig = &Cloud{
			name:  awsName,
			kind: AWS,
			Credentials: &aws.Credentials{
				AccessKeyId:     viper.GetString("AWS.AccessKeyId"),
				SecretAccessKey: viper.GetString("AWS.SecretAccessKey"),
			},
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloud)
		return
	}

	return
}

func (c *Cloud) Name() (cloudName string) {
	return c.name
}

func (c *Cloud) Type() (cloudType Kind) {
	return c.kind
}
