package workflows

import (
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type PackerTemplateFiles struct {
	General        *GeneralTemplateFile
	Aws            *AwsTemplateFile
	MergedTemplate *MergedTemplateFile
}

type PackerMergedTemplate struct {
	Instance    *template.Template
	Environment any
}

func NewPackerMergedTemplate(cloud binaries.Cloud) (packerMergedTemplate *PackerMergedTemplate, err error) {

	var (
		oopsBuilder = oops.
				Code("new_merged_template_failed").
				With("cloud", cloud)
		logger, _ = zap.NewProduction()

		packerGeneralEnvironment *PackerGeneralEnvironment
		packerTemplateFiles      *PackerTemplateFiles
		packerTemplateInstance   *template.Template
		absPathToTemplatesDir    string
	)

	defer logger.Sync()

	if absPathToTemplatesDir, err = filepath.Abs(binaries.TEMPLATE_DIR_NAME); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to create absolute path to: %s", binaries.TEMPLATE_DIR_NAME)
		return
	}

	packerTemplateFiles = &PackerTemplateFiles{
		General: &GeneralTemplateFile{
			Name:    PACKER_GENERAL_TEMPLATE_NAME,
			AbsPath: filepath.Join(absPathToTemplatesDir, PACKER_SUBDIR_NAME, binaries.GENERAL_SUBDIR_NAME, PACKER_GENERAL_TEMPLATE_NAME),
		},
		Aws: &AwsTemplateFile{
			Name:    PACKER_AWS_TEMPLATE_NAME,
			AbsPath: filepath.Join(absPathToTemplatesDir, PACKER_SUBDIR_NAME, binaries.AWS_SUBDIR_NAME, PACKER_AWS_TEMPLATE_NAME),
		},
		MergedTemplate: &MergedTemplateFile{
			Name:    PACKER_MERGED_TEMPLATE_NAME,
			AbsPath: filepath.Join(absPathToTemplatesDir, PACKER_SUBDIR_NAME, PACKER_MERGED_TEMPLATE_NAME),
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
			packerTemplateFiles.MergedTemplate.AbsPath,
			packerTemplateFiles.General.AbsPath,
			packerTemplateFiles.Aws.AbsPath,
		); err != nil {
			err = oopsBuilder.
				With("packerTemplates.General.AbsPath", packerTemplateFiles.General.AbsPath).
				With("packerTemplates.Aws.AbsPath", packerTemplateFiles.Aws.AbsPath).
				Wrapf(err, "failed to merge files to: %s", packerTemplateFiles.MergedTemplate.AbsPath)
			return
		}

		if packerTemplateInstance, err = template.ParseFiles(packerTemplateFiles.MergedTemplate.AbsPath); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while parsing template %s", packerTemplateFiles.MergedTemplate.AbsPath)
			return
		}

		packerMergedTemplate = &PackerMergedTemplate{
			Instance: packerTemplateInstance,
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

	if os.RemoveAll(pmt.Instance.Name()); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", pmt.Instance.Name())
		return
	}

	return
}

func (pmt *PackerMergedTemplate) Execute(writer io.Writer) (err error) {
	var (
		oopsBuilder = oops.
			Code("template_execute_failed").
			With("writer", writer)
	)

	if err = pmt.Instance.Execute(writer, pmt.Environment); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while executing template: %s", pmt.Instance.Name())
		return
	}

	return
}
