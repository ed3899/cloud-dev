package cloud

import "github.com/samber/oops"

type Credentials interface {
	Set() error
	Unset() error
}

type CloudDeployment struct {
	Kind        Kind
	Credentials Credentials
	Tool        Tool
}

func (cd *CloudDeployment) SetRunDir() (err error) {
	var (
		oopsBuilder = oops.Code(
			"cloud_deployment_set_run_dir_failed",
		)
	)

	switch cd.Kind {
	case AWS:
	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cd.Kind)
		return
	}
	return
}

func (cd *CloudDeployment) UnsetRunDir() (err error) {
	return
}
