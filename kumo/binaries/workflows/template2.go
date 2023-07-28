package workflows

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type GeneralTemplate struct {
	Name    string
	AbsPath string
}

type AwsTemplate struct {
	Name    string
	AbsPath string
}

type MergedTemplate struct {
	Name    string
	AbsPath string
}

type Template2 struct {
	Tool        binaries.Tool
	Cloud       binaries.Cloud
	AbsPath     string
	Environment any
}

func NewTemplate2(tool binaries.Tool, cloud binaries.Cloud) (mergedTemplate *Template2, err error) {
	const (
		// Packer
		packerSubDirName          = "packer"
		packerManifestName        = "manifest.json"
		packerGeneralTemplateName = "GeneralPackerVars.tmpl"
		packerAwsTemplateName     = "AWS_PackerVars.tmpl"
		packerMergedTemplateName  = "temp_merged_packer_template"
		// Terraform
		terraformSubDirName          = "terraform"
		terraformGeneralTemplateName = "GeneralTerraformTfVars.tmpl"
		terraformAwsTemplateName     = "AWS_TerraformTfVars.tmpl"
		terraformMergedTemplateName  = "temp_merged_terraform_template"
		// Subdirectory names
		generalSubDirName = "general"
		awsSubDirName     = "aws"
		templateDirName   = "templates"
		// Default IP
		defaulIp = "0.0.0.0"
	)

	var (
		oopsBuilder = oops.
				Code("new_merged_template_failed").
				With("tool", tool).
				With("cloud", cloud)
		logger, _ = zap.NewProduction()

		packerGeneralEnvironment    *PackerGeneralEnvironment
		packerManifest              *PackerManifest
		packerTemplates             *PackerTemplates
		terraformGeneralEnvironment *TerraformGeneralEnvironment
		terraformTemplates          *TerraformTemplates
		absPathToTemplatesDir       string
		publicIp                    string
		pickedIp                    string
		lastBuiltAmiId              string
		pickedAmiId                 string
	)

	defer logger.Sync()

	// Template paths
	if absPathToTemplatesDir, err = filepath.Abs(templateDirName); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to create absolute path to: %s", templateDirName)
		return
	}

	packerManifest = &PackerManifest{
		Name:    packerManifestName,
		AbsPath: filepath.Join(absPathToTemplatesDir, packerSubDirName, packerManifestName),
	}

	packerTemplates = &PackerTemplates{
		General: &GeneralTemplate{
			Name:    packerGeneralTemplateName,
			AbsPath: filepath.Join(absPathToTemplatesDir, packerSubDirName, generalSubDirName, packerGeneralTemplateName),
		},
		Aws: &AwsTemplate{
			Name:    packerAwsTemplateName,
			AbsPath: filepath.Join(absPathToTemplatesDir, packerSubDirName, awsSubDirName, packerAwsTemplateName),
		},
		MergedTemplate: &MergedTemplate{
			Name:    packerMergedTemplateName,
			AbsPath: filepath.Join(absPathToTemplatesDir, packerSubDirName, packerMergedTemplateName),
		},
	}

	terraformTemplates = &TerraformTemplates{
		General: &GeneralTemplate{
			Name:    terraformGeneralTemplateName,
			AbsPath: filepath.Join(absPathToTemplatesDir, terraformSubDirName, generalSubDirName, terraformGeneralTemplateName),
		},
		Aws: &AwsTemplate{
			Name:    terraformAwsTemplateName,
			AbsPath: filepath.Join(absPathToTemplatesDir, terraformSubDirName, awsSubDirName, terraformAwsTemplateName),
		},
		MergedTemplate: &MergedTemplate{
			Name:    terraformMergedTemplateName,
			AbsPath: filepath.Join(absPathToTemplatesDir, terraformSubDirName, terraformMergedTemplateName),
		},
	}

	switch tool {
	case binaries.PackerID:
		packerGeneralEnvironment = &PackerGeneralEnvironment{
			GIT_USERNAME:                          viper.GetString("Git.Username"),
			GIT_EMAIL:                             viper.GetString("Git.Email"),
			ANSIBLE_TAGS:                          viper.GetStringSlice("AMI.Tools"),
			GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
		}

		switch cloud {
		case binaries.AWS:
			// Merge templates
			if err = utils.MergeFilesTo(
				packerTemplates.MergedTemplate.AbsPath,
				packerTemplates.General.AbsPath,
				packerTemplates.Aws.AbsPath,
			); err != nil {
				err = oopsBuilder.
					With("packerTemplates.General.AbsPath", packerTemplates.General.AbsPath).
					With("packerTemplates.Aws.AbsPath", packerTemplates.Aws.AbsPath).
					Wrapf(err, "failed to merge files to: %s", packerTemplates.MergedTemplate.AbsPath)
				return
			}

			mergedTemplate = &Template2{
				Tool:    binaries.PackerID,
				Cloud:   binaries.AWS,
				AbsPath: packerTemplates.MergedTemplate.AbsPath,
				Environment: &PackerAWSEnvironment{
					AWS_ACCESS_KEY:                     viper.GetString("AWS.AccessKeyId"),
					AWS_SECRET_KEY:                     viper.GetString("AWS.SecretAccessKey"),
					AWS_IAM_PROFILE:                    viper.GetString("AWS.IamProfile"),
					AWS_USER_IDS:                       viper.GetStringSlice("AWS.UserIds"),
					AWS_AMI_NAME:                       viper.GetString("AMI.Name"),
					AWS_INSTANCE_TYPE:                  viper.GetString("AWS.EC2.Instance.Type"),
					AWS_REGION:                         viper.GetString("AWS.Region"),
					AWS_EC2_AMI_NAME_FILTER:            viper.GetString("AMI.Base.Filter"),
					AWS_EC2_AMI_ROOT_DEVICE_TYPE:       viper.GetString("AMI.Base.RootDeviceType"),
					AWS_EC2_AMI_VIRTUALIZATION_TYPE:    viper.GetString("AMI.Base.VirtualizationType"),
					AWS_EC2_AMI_OWNERS:                 viper.GetStringSlice("AMI.Base.Owners"),
					AWS_EC2_SSH_USERNAME:               viper.GetString("AMI.Base.User"),
					AWS_EC2_INSTANCE_USERNAME:          viper.GetString("AMI.User"),
					AWS_EC2_INSTANCE_USERNAME_HOME:     viper.GetString("AMI.Home"),
					AWS_EC2_INSTANCE_USERNAME_PASSWORD: viper.GetString("AMI.Password"),
					PackerGeneralEnvironment:           packerGeneralEnvironment,
				},
			}

		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud '%v' not supported", cloud)
			return
		}

	case binaries.TerraformID:
		// Get public IP
		if publicIp, err = utils.GetPublicIp(); err != nil {
			logger.Sugar().Warnf("failed to get public IP, using default: %s", defaulIp)
			pickedIp = defaulIp
		} else {
			pickedIp = publicIp
		}

		// Set general environment
		terraformGeneralEnvironment = &TerraformGeneralEnvironment{
			ALLOWED_IP: utils.MaskIp(pickedIp, 32),
		}

		switch cloud {
		case binaries.AWS:
			// Get last built AMI ID
			if lastBuiltAmiId, err = utils.GetLastBuiltAmiId(packerManifest.AbsPath); err != nil {
				err = oopsBuilder.
					Wrapf(err, "failed to get last built AMI ID")
				return
			}

			// Pick AMI ID to be used
			if pickedAmiId, err = utils.PickAmiIdToBeUsed(lastBuiltAmiId, viper.GetString("Up.AMI_Id")); err != nil {
				err = oopsBuilder.
					Wrapf(err, "failed to get picked AMI ID")
				return
			}

			// Merge templates
			if err = utils.MergeFilesTo(
				terraformTemplates.MergedTemplate.AbsPath,
				terraformTemplates.General.AbsPath,
				terraformTemplates.Aws.AbsPath,
			); err != nil {
				err = oopsBuilder.
					With("terraformTemplates.General.AbsPath", terraformTemplates.General.AbsPath).
					With("terraformTemplates.Aws.AbsPath", terraformTemplates.Aws.AbsPath).
					Wrapf(err, "failed to merge files to: %s", terraformTemplates.MergedTemplate.AbsPath)
				return
			}

			mergedTemplate = &Template2{
				Tool:    binaries.TerraformID,
				Cloud:   binaries.AWS,
				AbsPath: terraformTemplates.MergedTemplate.AbsPath,
				Environment: &TerraformAWSEnvironment{
					AWS_REGION:                   viper.GetString("AWS.Region"),
					AWS_INSTANCE_TYPE:            viper.GetString("AWS.EC2.Instance.Type"),
					AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
					AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
					AMI_ID:                       pickedAmiId,
					TerraformGeneralEnvironment:  terraformGeneralEnvironment,
				},
			}

		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud '%v' not supported", cloud)
			return
		}

	default:
		err = oopsBuilder.
			Wrapf(err, "Tool '%v' not supported", tool)
		return
	}

	return
}

func (t *Template2) Remove() (err error) {
	var (
		oopsBuilder = oops.
			Code("template_remove_failed")
	)

	if os.RemoveAll(t.AbsPath); err != nil {
		err = oopsBuilder.
			With("t.AbsPath", t.AbsPath).
			Wrapf(err, "Error occurred while removing %s", t.AbsPath)
		return
	}

	return
}
