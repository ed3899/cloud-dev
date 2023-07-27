package workflows

import (
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
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

type Template2 struct {
	AbsPath     string
	Environment any
}

type GeneralTemplate struct {
	Name    string
	AbsPath string
}

type AwsTemplate struct {
	Name    string
	AbsPath string
}

type MergedTemplate struct {
	Name		string
	AbsPath string
}

type PackerTemplates struct {
	General *GeneralTemplate
	Aws     *AwsTemplate
	MergedTemplate *MergedTemplate
}

type TerraformTemplates struct {
	General *GeneralTemplate
	Aws     *AwsTemplate
	MergedTemplate *MergedTemplate
}

func NewTemplate2(tool binaries.Tool, cloud binaries.Cloud) (mergedTemplate *Template2, err error) {
	const (
		packerSubDirName          = "packer"
		packerGeneralTemplateName = "GeneralPackerVars.tmpl"

		// Cloud
		packerAwsTemplateName    = "AWS_PackerVars.tmpl"
		packerMergedTemplateName = "temp_merged_packer_template"

		terraformSubDirName          = "terraform"
		terraformGeneralTemplateName = "GeneralTerraformTfVars.tmpl"
		terraformAwsTemplateName     = "AWS_TerraformTfVars.tmpl"
		terraformMergedTemplateName  = "temp_merged_terraform_template"

		generalSubDirName = "general"
		awsSubDirName     = "aws"

		templateDirName = "templates"
	)

	var (
		// absPathToPackerMergedTemplate    = filepath.Join(absPathToTemplatesDir, packerSubDirName, packerMergedTemplateName)
		// absPathToTerraformMergedTemplate = filepath.Join(absPathToTemplatesDir, terraformSubDirName, terraformMergedTemplateName)
		oopsBuilder                      = oops.
							Code("new_merged_template_failed").
							With("tool", tool).
							With("cloud", cloud)

		absPathToTemplatesDir string
	)

	// Template paths
	if absPathToTemplatesDir, err = filepath.Abs(templateDirName); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to create absolute path to: %s", templateDirName)
		return
	}

	packerTemplates := &PackerTemplates{
		General: &GeneralTemplate{
			Name:    packerGeneralTemplateName,
			AbsPath: filepath.Join(absPathToTemplatesDir, packerSubDirName, packerGeneralTemplateName),
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

	terraformTemplates := &TerraformTemplates{
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
		switch cloud {
		case binaries.AWS:
			if err = utils.MergeFilesTo(
				absPathToPackerMergedTemplate,
				packerTemplates.General.AbsPath,
				packerTemplates.Aws.AbsPath,
			); err != nil {
				err = oopsBuilder.
					With("packerTemplates.General.AbsPath", packerTemplates.General.AbsPath).
					With("packerTemplates.Aws.AbsPath", packerTemplates.Aws.AbsPath).
					Wrapf(err, "failed to merge files to: %s", absPathToPackerMergedTemplate)
				return
			}

			mergedTemplate = &MergedTemplate{
				AbsPath: mergedTemplateAbsPath,
			}

		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud '%v' not supported", cloud)
			return
		}
	case binaries.TerraformID:
		switch cloud {
		case binaries.AWS:
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

func (mt *MergedTemplate) Create() (err error) {
	return
}

func (mt *MergedTemplate) Remove() (err error) {
	return
}
