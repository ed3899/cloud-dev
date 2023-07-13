package templates

import (
	templates_terraform_generic "github.com/ed3899/kumo/templates/terraform/generic"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type AWS_TerraformEnvironment struct {
	AWS_REGION                   string
	AWS_INSTANCE_TYPE            string
	AWS_EC2_INSTANCE_VOLUME_TYPE string
	AWS_EC2_INSTANCE_VOLUME_SIZE int
}

func CraftAWS_TerraformTfVarsFile(cloud string) (generalTerraformVarsPath string, err error) {
	awsEnv := &AWS_TerraformEnvironment{
		AWS_REGION:                   viper.GetString("AWS.Region"),
		AWS_INSTANCE_TYPE:            viper.GetString("AWS.EC2.Instance.Type"),
		AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
		AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
	}

	awsTerraformVarsPath, err := templates_terraform_generic.CraftGenericCloudTerraformTfVarsFile(cloud, "AWS_TerraformTfVars.tmpl", "aws.auto.tfvars", *awsEnv)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting generic Terraform Vars file for cloud '%s'", cloud)
		return "", err
	}

	return awsTerraformVarsPath, nil
}
