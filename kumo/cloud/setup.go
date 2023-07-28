package cloud

import (
	"github.com/ed3899/kumo/cloud/aws"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Credentials interface {
	Set() error
	Unset() error
}

type CloudSetup struct {
	cloudName string

	Credentials Credentials
}

func (cs *CloudSetup) GetCloudName() (cloudName string) {
	return cs.cloudName
}

func NewCloudSetup(cloud string) (cloudSetup *CloudSetup, err error) {
	var (
		oopsBuilder = oops.
				Code("new_cloud_deployment_failed").
				With("cloud", cloud)

		cloudName   string
		credentials Credentials
	)

	switch cloud {
	case "aws":
		cloudName = "aws"
		credentials = &aws.Credentials{
			AccessKeyId:     viper.GetString("AWS.AccessKeyId"),
			SecretAccessKey: viper.GetString("AWS.SecretAccessKey"),
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloud)
		return
	}

	cloudSetup = &CloudSetup{
		cloudName: cloudName,

		Credentials: credentials,
	}

	return
}
