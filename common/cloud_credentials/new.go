package cloud_credentials

import (
	"github.com/ed3899/kumo/common/cloud_config"
	"github.com/ed3899/kumo/common/cloud_credentials/aws"
	"github.com/ed3899/kumo/common/cloud_credentials/interfaces"
	"github.com/samber/oops"
)

func New(cloud cloud_config.CloudI) (credentials interfaces.Credentials, err error) {
	var (
		oopsBuilder = oops.
			Code("new_credentials_failed").
			With("cloud", cloud)
	)

	switch cloud.Type() {
	case cloud_config.AWS:
		credentials = aws.New(cloud)
	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloud.Name())
		return
	}

	return
}
