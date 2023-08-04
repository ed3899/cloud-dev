package cloud_credentials

import (
	common_cloud_constants "github.com/ed3899/kumo/common/cloud/constants"
	common_cloud_interfaces "github.com/ed3899/kumo/common/cloud/interfaces"
	"github.com/ed3899/kumo/common/cloud_credentials/aws"
	"github.com/ed3899/kumo/common/cloud_credentials/interfaces"
	"github.com/samber/oops"
)

func New(cloud common_cloud_interfaces.Cloud) (credentials interfaces.Credentials, err error) {
	var (
		oopsBuilder = oops.
			Code("new_credentials_failed").
			With("cloud", cloud)
	)

	switch cloud.Type() {
	case common_cloud_constants.AWS:
		credentials = aws.New(cloud)
	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloud.Name())
		return
	}

	return
}
