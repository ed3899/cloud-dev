package environment

import (
	"github.com/ed3899/kumo/tool/environment/packer/general"
)

type CloudEnvironmentI interface {
	IsCloudEnvironment() (isCloudEnvironment bool)
}

type Environment struct {
	General general.PackerGeneralEnvironment
	Cloud   CloudEnvironmentI
}

func NewEnvironment(
	newGeneralEnvironment general.NewEnvironmentF,
	cloudEnvironment CloudEnvironmentI,
) (
	environment Environment,
) {

	environment = Environment{
		General: newGeneralEnvironment(),
		Cloud:   cloudEnvironment,
	}

	return

}
