package environment

import (
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/utils/ip"
	"go.uber.org/zap"
)

func NewEnvironment(
	getPublicIp ip.GetPublicIpF,
	withMask ip.WithMaskF,
) (
	generalEnvironment Environment,
) {
	var (
		logger, _ = zap.NewProduction()

		err      error
		publicIp string
		pickedIp string
	)

	defer logger.Sync()

	if publicIp, err = getPublicIp(); err != nil {
		logger.Warn(
			"There was an error getting the public ip, using default instead",
			zap.String("error", err.Error()),
			zap.String("defaultIp", constants.TERRAFORM_DEFAULT_ALLOWED_IP),
		)
		pickedIp = constants.TERRAFORM_DEFAULT_ALLOWED_IP
	} else {
		pickedIp = publicIp
	}

	generalEnvironment = Environment{
		Required: Required{
			ALLOWED_IP: withMask(32)(pickedIp),
		},
	}

	return
}

type Required struct {
	ALLOWED_IP string
}

type Environment struct {
	Required Required
}

type NewTerraformGeneralEnvironmentF func(GetPublicIp ip.GetPublicIpF, WithMask ip.WithMaskF) (terraformGeneralEnvironment Environment)
