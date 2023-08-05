package environment

import (
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/utils"
	"go.uber.org/zap"
)

func NewTerraformGeneralEnvironment(
	GetPublicIp utils.GetPublicIpF,
	WithMask utils.WithMaskF,
) (
	terraformGeneralEnvironment TerraformGeneralEnvironment,
) {
	var (
		logger, _ = zap.NewProduction()
		
		err       error
		publicIp       string
		pickedIp       string
		MaskIpWith32   utils.MaskIpL
		maskedIpWith32 string
	)

	defer logger.Sync()

	if publicIp, err = GetPublicIp(); err != nil {
		logger.Warn(
			"There was an error getting the public ip, using default instead",
			zap.String("error", err.Error()),
			zap.String("defaultIp", constants.TERRAFORM_DEFAULT_ALLOWED_IP),
		)
		pickedIp = constants.TERRAFORM_DEFAULT_ALLOWED_IP
	} else {
		pickedIp = publicIp
	}

	MaskIpWith32 = WithMask(32)
	maskedIpWith32 = MaskIpWith32(pickedIp)

	terraformGeneralEnvironment = TerraformGeneralEnvironment{
		Required: TerraformGeneralRequired{
			ALLOWED_IP: maskedIpWith32,
		},
	}

	return
}
