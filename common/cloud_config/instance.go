package cloud_config

import (
	"github.com/samber/oops"
)

type Cloud struct {
	name string
	kind Kind
}

func New(cloudFromConfig string) (cloud *Cloud, err error) {
	var (
		oopsBuilder = oops.
				Code("new_cloud_deployment_failed").
				With("cloud", cloudFromConfig)
		awsName = AWS_NAME
	)

	switch cloudFromConfig {
	case awsName:
		cloud = &Cloud{
			name: awsName,
			kind: AWS,
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

func (c *Cloud) Type() (cloudType Kind) {
	return c.kind
}
