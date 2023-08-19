package environment

import (
	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/utils/ip"
	"go.uber.org/zap"
)

type TerraformBaseRequired struct {
	ALLOWED_IP string
}

type TerraformBaseEnvironment struct {
	Required *TerraformBaseRequired
}

func NewTerraformBaseEnvironment() *TerraformBaseEnvironment {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var pickedIp string

	publicIp, err := ip.GetPublicIp()
	if err != nil {
		logger.Warn("Failed to get public IP address", zap.Error(err))
		logger.Warn("Using default ip instead", zap.String("ip", constants.TERRAFORM_DEFAULT_ALLOWED_IP))

		pickedIp = constants.TERRAFORM_DEFAULT_ALLOWED_IP
	} else {
		pickedIp = publicIp
	}

	return &TerraformBaseEnvironment{
		Required: &TerraformBaseRequired{
			ALLOWED_IP: ip.MaskIp(pickedIp, 32),
		},
	}
}
