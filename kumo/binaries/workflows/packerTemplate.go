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

type PackerTemplates struct {
	General        *GeneralTemplate
	Aws            *AwsTemplate
	MergedTemplate *MergedTemplate
}

type PackerMergedTemplate struct {
	Name        string
	AbsPath     string
	Environment any
}

func NewPackerMergedTemplate(cloud binaries.Cloud) (packerMergedTemplate *PackerMergedTemplate, err error) {
	const (
		subDirName          = "packer"
		manifestName        = "manifest.json"
		generalTemplateName = "GeneralPackerVars.tmpl"
		awsTemplateName     = "AWS_PackerVars.tmpl"
		mergedTemplateName  = "temp_merged_packer_template"
	)

	var (
		oopsBuilder = oops.
				Code("new_merged_template_failed").
				With("cloud", cloud)
		logger, _ = zap.NewProduction()

		packerGeneralEnvironment *PackerGeneralEnvironment
		packerTemplates          *PackerTemplates
		absPathToTemplatesDir    string
	)

	defer logger.Sync()

	if absPathToTemplatesDir, err = filepath.Abs(binaries.TEMPLATE_DIR_NAME); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to create absolute path to: %s", binaries.TEMPLATE_DIR_NAME)
		return
	}

	packerTemplates = &PackerTemplates{
		General: &GeneralTemplate{
			Name:    generalTemplateName,
			AbsPath: filepath.Join(absPathToTemplatesDir, subDirName, binaries.GENERAL_SUBDIR_NAME, generalTemplateName),
		},
		Aws: &AwsTemplate{
			Name:    awsTemplateName,
			AbsPath: filepath.Join(absPathToTemplatesDir, subDirName, binaries.GENERAL_SUBDIR_NAME, generalTemplateName),
		},
		MergedTemplate: &MergedTemplate{
			Name:    mergedTemplateName,
			AbsPath: filepath.Join(absPathToTemplatesDir, subDirName, mergedTemplateName),
		},
	}

	packerGeneralEnvironment = &PackerGeneralEnvironment{
		GIT_USERNAME:                          viper.GetString("Git.Username"),
		GIT_EMAIL:                             viper.GetString("Git.Email"),
		ANSIBLE_TAGS:                          viper.GetStringSlice("AMI.Tools"),
		GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
	}

	switch cloud {
	case binaries.AWS:
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

		packerMergedTemplate = &PackerMergedTemplate{
			Name:    mergedTemplateName,
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
			Errorf("Cloud '%v' not supported", cloud)
		return
	}

	return
}

func (pmt *PackerMergedTemplate) Remove() (err error) {
	var (
		oopsBuilder = oops.
			Code("template_remove_failed")
	)

	if os.RemoveAll(pmt.AbsPath); err != nil {
		err = oopsBuilder.
			With("t.AbsPath", pmt.AbsPath).
			Wrapf(err, "Error occurred while removing %s", pmt.AbsPath)
		return
	}

	return
}
