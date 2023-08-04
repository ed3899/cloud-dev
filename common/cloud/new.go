package cloud

import (
	"github.com/ed3899/kumo/common/cloud/constants"
	"github.com/samber/oops"
)

type Cloud struct {
	name string
	kind constants.Kind
}

func New(cloudFromConfig string) (cloud *Cloud, err error) {
	var (
		oopsBuilder = oops.
				Code("new_cloud_deployment_failed").
				With("cloud", cloudFromConfig)
		awsName = constants.AWS_NAME
	)

	switch cloudFromConfig {
	case awsName:
		cloud = &Cloud{
			name: awsName,
			kind: constants.AWS,
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloudFromConfig)
		return
	}

	return
}

func (c *Cloud) Name() (cloudName string) {
	return c.name
}

func (c *Cloud) Kind() (cloudKind constants.Kind) {
	return c.kind
}
