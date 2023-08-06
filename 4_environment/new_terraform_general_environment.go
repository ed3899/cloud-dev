package environment

import (
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/utils"
	"go.uber.org/zap"
)

func NewTerraformGeneralEnvironment(
	getPublicIp utils.GetPublicIpF,
	withMask utils.WithMaskF,
) (
	terraformGeneralEnvironment TerraformGeneralEnvironmentI,
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

	terraformGeneralEnvironment = TerraformGeneralEnvironment{
		Required: TerraformGeneralRequired{
			ALLOWED_IP: withMask(32)(pickedIp),
		},
	}

	return
}

type TerraformGeneralRequired struct {
	ALLOWED_IP string
}

type TerraformGeneralEnvironment struct {
	Required TerraformGeneralRequired
}

type TerraformGeneralEnvironmentI interface {
	IsTerraformGeneralEnvironment() bool
}

func (tae TerraformGeneralEnvironment) IsTerraformGeneralEnvironment() bool {
	return true
}

type NewTerraformGeneralEnvironmentF func(GetPublicIp utils.GetPublicIpF, WithMask utils.WithMaskF) (terraformGeneralEnvironment TerraformGeneralEnvironmentI)
