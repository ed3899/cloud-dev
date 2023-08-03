package cloud

import (
	"github.com/ed3899/kumo/common/cloud/aws"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Credentials interface {
	Set() error
	Unset() error
}

type Cloud struct {
	name   string
	type_   Type
	Credentials Credentials
}

func (c *Cloud) GetName() (cloudName string) {
	return c.name
}

func (c *Cloud) Type() (cloudType Type) {
	return c.type_
}

func NewConfig(cloud string) (cloudConfig *Cloud, err error) {
	var (
		oopsBuilder = oops.
				Code("new_cloud_deployment_failed").
				With("cloud", cloud)
		awsName = AWS_NAME
	)

	switch cloud {
	case awsName:
		cloudConfig = &Cloud{
			name: awsName,
			type_: AWS,
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
