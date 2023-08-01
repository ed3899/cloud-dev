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

type CloudSetup struct {
	cloudName   string
	cloudType   CloudType
	Credentials Credentials
}

func (cs *CloudSetup) GetCloudName() (cloudName string) {
	return cs.cloudName
}

func (cs *CloudSetup) GetCloudType() (cloudType CloudType) {
	return cs.cloudType
}

func NewCloudSetup(cloud string) (cloudSetup *CloudSetup, err error) {
	var (
		oopsBuilder = oops.
			Code("new_cloud_deployment_failed").
			With("cloud", cloud)
	)

	switch cloud {
	case "aws":
		cloudSetup = &CloudSetup{
			cloudName: "aws",
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
