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

type Config struct {
	cloudName   string
	cloudType   CloudType
	Credentials Credentials
}

func (c *Config) GetCloudName() (cloudName string) {
	return c.cloudName
}

func (c *Config) GetCloudType() (cloudType CloudType) {
	return c.cloudType
}

func NewConfig(cloud string) (cloudConfig *Config, err error) {
	var (
		oopsBuilder = oops.
				Code("new_cloud_deployment_failed").
				With("cloud", cloud)
		awsName = AWS_NAME
	)

	switch cloud {
	case awsName:
		cloudConfig = &Config{
			cloudName: awsName,
			cloudType: AWS,
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
