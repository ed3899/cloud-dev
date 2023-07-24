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

type Tool int

const (
	PackerID Tool = iota
	TerraformID
)

type Cloud int

const (
	AWS Cloud = iota
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
	Instance    *template.Template
	Environment any
}

type HashicorpVars struct {
	Name     string
	AbsPath  string
	Template *Template
}

func NewHashicorpVars(tool Tool, cloud Cloud) (vf *HashicorpVars, err error) {
	const (
		packerSubDirName          = "packer"
		packerGeneralVarsFileName = "general_ami.auto.pkrvars.hcl"
		packerGeneralTemplateName = "GeneralPackerVars.tmpl"
		packerAwsVarsFileName     = "aws_ami.auto.pkrvars.hcl"
		packerAwsTemplateName     = "AWS_PackerVars.tmpl"
		packerManifestName        = "manifest.json"
	)
	const (
		terraformSubDirName          = "terraform"
		terraformGeneralVarsFileName = "general.auto.tfvars"
		terraformGeneralTemplateName = "GeneralTerraformTfVars.tmpl"
		terraformAwsVarsFileName     = "aws.auto.tfvars"
		terraformAwsTemplateName     = "AWS_TerraformTfVars.tmpl"
	)
	const (
		generalSubDirName = "general"
		awsSubDirName     = "aws"
	)

	const (
		defaulIp = "0.0.0.0"
	)

	var (
		packerGeneralEnvironment    *PackerGeneralEnvironment
		terraformGeneralEnvironment *TerraformGeneralEnvironment
		mergedTemplateAbsPath       string
	)

	var (
		pickedIp    string
		pickedAmiId string
	)

	absPathToTemplatesDir, err := filepath.Abs(filepath.Join("templates"))
	if err != nil {
		err = errors.Wrap(err, "failed to create path to templates")
		return nil, err
	}

	var (
		absPathToPackerGeneralTemplate    = filepath.Join(absPathToTemplatesDir, packerSubDirName, packerGeneralTemplateName)
		absPathToPackerAWSTemplate        = filepath.Join(absPathToTemplatesDir, packerSubDirName, packerAwsTemplateName)
		absPathToTerraformGeneralTemplate = filepath.Join(absPathToTemplatesDir, terraformSubDirName, terraformGeneralTemplateName)
		absPathToTerraformAWSTemplate     = filepath.Join(absPathToTemplatesDir, terraformSubDirName, terraformAwsTemplateName)
	)

	absPathToPackerDir, err := filepath.Abs(filepath.Join(packerSubDirName))
	if err != nil {
		err = errors.Wrap(err, "failed to create path to packer directory")
		return nil, err
	}

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
			mergedTemplateAbsPath, err = utils.MergeFiles(absPathToPackerGeneralTemplate, absPathToPackerAWSTemplate)
			if err != nil {
				err = errors.Wrap(err, "failed to merge files")
				return nil, err
			}

			vf.Name = packerAwsVarsFileName
			vf.AbsPath = mergedTemplateAbsPath
			vf.Template = &Template{
				AbsPath: absPathToPackerAWSTemplate,
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
			err = errors.Errorf("Kind '%v' not supported", cloud)
			return
		}
	case TerraformID:
		publicIp, err := utils.GetPublicIp()
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
			mergedTemplateAbsPath, err = utils.MergeFiles(absPathToTerraformGeneralTemplate, absPathToTerraformAWSTemplate)
			if err != nil {
				err = errors.Wrap(err, "failed to merge files")
				return nil, err
			}

			absPathToPackerManifestForAws := filepath.Join(absPathToPackerDir, awsSubDirName, packerManifestName)

			lastBuiltAmiId, err := utils.GetLastBuiltAmiId(absPathToPackerManifestForAws)
			if err != nil {
				err = errors.Wrap(err, "failed to get last built AMI ID")
				return nil, err
			}

			pickedAmiId, err = utils.PickAmiIdToBeUsed(lastBuiltAmiId, utils.GetAmiIdFromConfig())
			if err != nil {
				err = errors.Wrap(err, "failed to pick AMI ID to be used")
				return nil, err
			}

			vf.Name = terraformAwsVarsFileName
			vf.AbsPath = mergedTemplateAbsPath
			vf.Template = &Template{
				AbsPath: absPathToTerraformAWSTemplate,
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
			err = errors.Errorf("Kind '%v' not supported", cloud)
			return nil, err
		}
	default:
		err = errors.Errorf("Tool '%s' not supported", tool)
		return
	}

	return
}

func (hv *HashicorpVars) CreateFile() (err error) {
	// Create vars file to be filled with values
	varsFile, err := os.Create(hv.AbsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while creating %s", hv.AbsPath)
		return
	}
	defer varsFile.Close()

	// Get template
	template, err := template.ParseFiles(hv.Template.AbsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s", hv.Template.AbsPath)
		return
	}

	// Execute template
	err = hv.Template.Instance.Execute(varsFile, template)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while executing template: %s", hv.Template.AbsPath)
		return
	}

	return
}
