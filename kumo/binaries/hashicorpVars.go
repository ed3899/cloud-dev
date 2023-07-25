package binaries

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type PackerGeneralEnvironment struct {
	GIT_USERNAME                          string
	GIT_EMAIL                             string
	ANSIBLE_TAGS                          []string
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

type PackerAWSEnvironment struct {
	AWS_ACCESS_KEY                     string
	AWS_SECRET_KEY                     string
	AWS_IAM_PROFILE                    string
	AWS_USER_IDS                       []string
	AWS_AMI_NAME                       string
	AWS_INSTANCE_TYPE                  string
	AWS_REGION                         string
	AWS_EC2_AMI_NAME_FILTER            string
	AWS_EC2_AMI_ROOT_DEVICE_TYPE       string
	AWS_EC2_AMI_VIRTUALIZATION_TYPE    string
	AWS_EC2_AMI_OWNERS                 []string
	AWS_EC2_SSH_USERNAME               string
	AWS_EC2_INSTANCE_USERNAME          string
	AWS_EC2_INSTANCE_USERNAME_HOME     string
	AWS_EC2_INSTANCE_USERNAME_PASSWORD string
	PackerGeneralEnvironment           *PackerGeneralEnvironment
}

type TerraformGeneralEnvironment struct {
	ALLOWED_IP string
}

type TerraformAWSEnvironment struct {
	AWS_REGION                   string
	AWS_INSTANCE_TYPE            string
	AWS_EC2_INSTANCE_VOLUME_TYPE string
	AWS_EC2_INSTANCE_VOLUME_SIZE int
	AMI_ID                       string
	TerraformGeneralEnvironment  *TerraformGeneralEnvironment
}

type Template struct {
	AbsPath     string
	Environment any
}

func (t *Template) Remove() (err error) {
	if os.RemoveAll(t.AbsPath); err != nil {
		return errors.Wrapf(err, "Error occurred while removing %s", t.AbsPath)
	}
	return
}

type HashicorpVars struct {
	Name     string
	AbsPath  string
	Template *Template
}

func NewHashicorpVars(tool Tool, cloud Cloud) (hv *HashicorpVars, err error) {
	const (
		//% Packer
		packerSubDirName   = "packer"
		packerManifestName = "manifest.json"
		// General
		packerGeneralVarsFileName = "general_ami.auto.pkrvars.hcl"
		packerGeneralTemplateName = "GeneralPackerVars.tmpl"
		// AWS
		packerAwsVarsFileName = "aws_ami.auto.pkrvars.hcl"
		packerAwsTemplateName = "AWS_PackerVars.tmpl"
		// Merged
		packerMergedTemplateName = "temp_merged_packer_template"

		//% Terraform
		terraformSubDirName = "terraform"
		// General
		terraformGeneralVarsFileName = "general.auto.tfvars"
		terraformGeneralTemplateName = "GeneralTerraformTfVars.tmpl"
		// AWS
		terraformAwsVarsFileName = "aws.auto.tfvars"
		terraformAwsTemplateName = "AWS_TerraformTfVars.tmpl"
		// Merged
		terraformMergedTemplateName = "temp_merged_terraform_template"

		// Subdirectory names
		generalSubDirName = "general"
		awsSubDirName     = "aws"

		defaulIp        = "0.0.0.0"
		templateDirName = "templates"
	)

	var (
		// Environments
		packerGeneralEnvironment    *PackerGeneralEnvironment
		terraformGeneralEnvironment *TerraformGeneralEnvironment

		// Choices
		publicIp       string
		pickedIp       string
		lastBuiltAmiId string
		pickedAmiId    string

		// Resulting merged template
		mergedTemplateAbsPath string

		// Run directories
		absPathToPackerRunDir    string
		absPathToTerraformRunDir string
	)

	if absPathToPackerRunDir, err = filepath.Abs(packerSubDirName); err != nil {
		err = errors.Wrap(err, "failed to create path to packer run directory")
		return nil, err
	}

	if absPathToTerraformRunDir, err = filepath.Abs(terraformSubDirName); err != nil {
		err = errors.Wrap(err, "failed to create path to terraform run directory")
		return nil, err
	}

	// Var file paths
	var (
		absPathToPackerAwsVarsFile = filepath.Join(absPathToPackerRunDir, awsSubDirName, packerAwsVarsFileName)

		absPathToTerraformAwsVarsFile = filepath.Join(absPathToTerraformRunDir, awsSubDirName, terraformAwsVarsFileName)
	)

	// Template paths
	absPathToTemplatesDir, err := filepath.Abs(templateDirName)
	if err != nil {
		err = errors.Wrap(err, "failed to create path to templates")
		return nil, err
	}

	var (
		absPathToPackerGeneralTemplate    = filepath.Join(absPathToTemplatesDir, packerSubDirName, generalSubDirName, packerGeneralTemplateName)
		absPathToPackerAWSTemplate        = filepath.Join(absPathToTemplatesDir, packerSubDirName, awsSubDirName, packerAwsTemplateName)
		absPathToTempPackerMergedTemplate = filepath.Join(absPathToTemplatesDir, packerSubDirName, packerMergedTemplateName)

		absPathToTerraformGeneralTemplate    = filepath.Join(absPathToTemplatesDir, terraformSubDirName, generalSubDirName, terraformGeneralTemplateName)
		absPathToTerraformAWSTemplate        = filepath.Join(absPathToTemplatesDir, terraformSubDirName, awsSubDirName, terraformAwsTemplateName)
		absPathToTempTerraformMergedTemplate = filepath.Join(absPathToTemplatesDir, terraformSubDirName, terraformMergedTemplateName)
	)

	// Packer manifest paths
	var (
		absPathToPackerAwsManifest = filepath.Join(absPathToPackerRunDir, awsSubDirName, packerManifestName)
	)

	switch tool {
	case PackerID:
		packerGeneralEnvironment = &PackerGeneralEnvironment{
			GIT_USERNAME:                          viper.GetString("Git.Username"),
			GIT_EMAIL:                             viper.GetString("Git.Email"),
			ANSIBLE_TAGS:                          viper.GetStringSlice("AMI.Tools"),
			GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
		}
		switch cloud {
		case AWS:
			// Create merged template
			mergedTemplateAbsPath, err = utils.MergeFilesTo(
				absPathToTempPackerMergedTemplate,
				absPathToPackerGeneralTemplate,
				absPathToPackerAWSTemplate,
			)
			if err != nil {
				err = errors.Wrap(err, "failed to merge files")
				return nil, err
			}

			hv = &HashicorpVars{
				Name:    packerGeneralVarsFileName,
				AbsPath: absPathToPackerAwsVarsFile,
				Template: &Template{
					AbsPath: mergedTemplateAbsPath,
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
				},
			}
		default:
			err = errors.Errorf("Kind '%v' not supported", cloud)
		}
	case TerraformID:
		publicIp, err = utils.GetPublicIp()
		if err != nil {
			err = errors.Wrapf(err, "failed to get public IP, using default: %s", defaulIp)
			log.Print(err)
			pickedIp = defaulIp
		}

		pickedIp = publicIp
		terraformGeneralEnvironment = &TerraformGeneralEnvironment{
			ALLOWED_IP: pickedIp,
		}

		switch cloud {
		case AWS:
			mergedTemplateAbsPath, err = utils.MergeFilesTo(
				absPathToTempTerraformMergedTemplate,
				absPathToTerraformGeneralTemplate,
				absPathToTerraformAWSTemplate,
			)
			if err != nil {
				err = errors.Wrap(err, "failed to merge files")
				return nil, err
			}

			lastBuiltAmiId, err = utils.GetLastBuiltAmiId(absPathToPackerAwsManifest)
			if err != nil {
				err = errors.Wrap(err, "failed to get last built AMI ID")
				return nil, err
			}

			pickedAmiId, err = utils.PickAmiIdToBeUsed(lastBuiltAmiId, utils.GetAmiIdFromConfig())
			if err != nil {
				err = errors.Wrap(err, "failed to pick AMI ID to be used")
				return nil, err
			}

			hv = &HashicorpVars{
				Name:    terraformAwsVarsFileName,
				AbsPath: absPathToTerraformAwsVarsFile,
				Template: &Template{
					AbsPath: mergedTemplateAbsPath,
					Environment: &TerraformAWSEnvironment{
						AWS_REGION:                   viper.GetString("AWS.Region"),
						AWS_INSTANCE_TYPE:            viper.GetString("AWS.EC2.Instance.Type"),
						AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
						AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
						AMI_ID:                       pickedAmiId,
						TerraformGeneralEnvironment:  terraformGeneralEnvironment,
					},
				},
			}
		default:
			err = errors.Errorf("Kind '%v' not supported", cloud)
		}
	default:
		err = errors.Errorf("Tool '%v' not supported", tool)
	}

	return
}

func (hv *HashicorpVars) Create() (err error) {
	// Create vars file to be filled with values
	varsFile, err := os.Create(hv.AbsPath)
	if err != nil {
		return errors.Wrapf(err, "Error occurred while creating %s", hv.AbsPath)
	}
	defer varsFile.Close()

	// Get template and defer deletion
	template, err := template.ParseFiles(hv.Template.AbsPath)
	if err != nil {
		return errors.Wrapf(err, "Error occurred while crafting absolute path to %s", hv.Template.AbsPath)
	}
	defer func() {
		if err = hv.Template.Remove(); err != nil {
			err = errors.Wrapf(err, "Error occurred while removing %s", hv.Template.AbsPath)
			return
		}
	}()

	// Execute template
	err = template.Execute(varsFile, hv.Template.Environment)
	if err != nil {
		return errors.Wrapf(err, "Error occurred while executing template: %s", hv.Template.AbsPath)
	}

	return
}
