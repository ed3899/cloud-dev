package templates

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

// Crafts a generic cloud Packer Vars file.
//
// It looks for templates under "./templates/packer/<cloud>/<templateName>.tmpl"
//
// It creates a file under "./packer/<cloud>/<packerVarsFileName>.auto.tfvars"
//
// The env parameter is a pointer to a struct that contains the data
// to be used in the template. Example:
//
//	awsEnv := &AWS_PackerEnvironment{
//		AWS_ACCESS_KEY:                     viper.GetString("AWS.AccessKeyId"),
//		AWS_SECRET_KEY:                     viper.GetString("AWS.SecretAccessKey"),
//		AWS_IAM_PROFILE:                    viper.GetString("AWS.IamProfile"),
//		AWS_USER_IDS:                       viper.GetStringSlice("AWS.UserIds"),
//		AWS_AMI_NAME:                       viper.GetString("AMI.Name"),
//		AWS_INSTANCE_TYPE:                  viper.GetString("AWS.EC2.Instance.Type"),
//		AWS_REGION:                         viper.GetString("AWS.Region"),
//		AWS_EC2_AMI_NAME_FILTER:            viper.GetString("AMI.Base.Filter"),
//		AWS_EC2_AMI_ROOT_DEVICE_TYPE:       viper.GetString("AMI.Base.RootDeviceType"),
//		AWS_EC2_AMI_VIRTUALIZATION_TYPE:    viper.GetString("AMI.Base.VirtualizationType"),
//		AWS_EC2_AMI_OWNERS:                 viper.GetStringSlice("AMI.Base.Owners"),
//		AWS_EC2_SSH_USERNAME:               viper.GetString("AMI.Base.User"),
//		AWS_EC2_INSTANCE_USERNAME:          viper.GetString("AMI.User"),
//		AWS_EC2_INSTANCE_USERNAME_HOME:     viper.GetString("AMI.Home"),
//		AWS_EC2_INSTANCE_USERNAME_PASSWORD: viper.GetString("AMI.Password"),
//	}
//
// Usage:
//
//	("aws", "AWS_PackerVarsTemplate.tmpl", "aws.auto.pkrvars", env) -> ("packer/aws/aws.auto.pkrvars", nil)
func CraftGenericCloudPackerVarsFile(cloud, templateName, packerVarsFileName string, env any) (resultingPackerVarsPath string, err error) {
	// Get template
	templatePath, err := filepath.Abs(filepath.Join("templates", "packer", cloud, templateName))
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s template file", templateName)
		return "", err
	}

	// Create vars file
	resultingPackerVarsPath, err = filepath.Abs(filepath.Join("packer", cloud, packerVarsFileName))
	if err != nil {
		err = errors.Wrap(err, "Error occurred while crafting absolute path to Packer AWS Vars file")
		return "", err
	}

	file, err := os.Create(resultingPackerVarsPath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while creating Packer AWS Vars file")
		return "", err
	}
	defer file.Close()

	// Execute template file
	err = tmpl.Execute(file, env)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while executing Packer AWS Vars template file")
		return "", err
	}

	// Return path to vars file
	return resultingPackerVarsPath, nil
}
