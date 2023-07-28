package workflows

import (
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/utils"
	"github.com/samber/oops"
	"go.uber.org/zap"
)

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

type TerraformTemplates struct {
	General        *GeneralTemplate
	Aws            *AwsTemplate
	MergedTemplate *MergedTemplate
}

type PackerManifest struct {
	Name    string
	AbsPath string
}

type TerraformMergedTemplate struct {
	Name        string
	AbsPath     string
	Environment any
}

func NewTerraformMergedTemplate(cloud binaries.Cloud) (terraformMergedTemplate *TerraformMergedTemplate, err error) {
	var (
		oopsBuilder = oops.
				Code("new_terraform_merged_template_failed")
		logger, _ = zap.NewProduction()

		terraformGeneralEnvironment *TerraformGeneralEnvironment
		terraformTemplates          *TerraformTemplates
		packerManifest              *PackerManifest
		absPathToTemplatesDir       string
		publicIp                    string
		pickedIp                    string
		lastBuiltAmiId              string
		pickedAmiId                 string
	)

	defer logger.Sync()

	if absPathToTemplatesDir, err = filepath.Abs(binaries.TEMPLATE_DIR_NAME); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to create absolute path to: %s", binaries.TEMPLATE_DIR_NAME)
		return
	}

	terraformTemplates = &TerraformTemplates{
		General: &GeneralTemplate{
			Name:    TERRAFORM_GENERAL_TEMPLATE_NAME,
			AbsPath: filepath.Join(absPathToTemplatesDir, TERRAFORM_SUBDIR_NAME, binaries.GENERAL_SUBDIR_NAME, TERRAFORM_GENERAL_TEMPLATE_NAME),
		},
		Aws: &AwsTemplate{
			Name:    TERRAFORM_AWS_TEMPLATE_NAME,
			AbsPath: filepath.Join(absPathToTemplatesDir, TERRAFORM_SUBDIR_NAME, binaries.AWS_SUBDIR_NAME, TERRAFORM_AWS_TEMPLATE_NAME),
		},
		MergedTemplate: &MergedTemplate{
			Name:    TERRAFORM_MERGED_TEMPLATE_NAME,
			AbsPath: filepath.Join(absPathToTemplatesDir, TERRAFORM_SUBDIR_NAME, TERRAFORM_MERGED_TEMPLATE_NAME),
		},
	}

	packerManifest = &PackerManifest{
		Name:    PACKER_MANIFEST_NAME,
		AbsPath: filepath.Join(absPathToTemplatesDir, PACKER_SUBDIR_NAME, PACKER_MANIFEST_NAME),
	}

	// Get public IP
	if publicIp, err = utils.GetPublicIp(); err != nil {
		logger.Sugar().Warnf("failed to get public IP, using default: %s", binaries.DEFAULT_IP)
		pickedIp = binaries.DEFAULT_IP
	} else {
		pickedIp = publicIp
	}

	// Set general environment
	terraformGeneralEnvironment = &TerraformGeneralEnvironment{
		ALLOWED_IP: utils.MaskIp(pickedIp, 32),
	}

	switch cloud {

	case binaries.AWS:
	default:
		err = oopsBuilder.
			Wrapf(err, "Cloud '%v' not supported", cloud)
		return
	}

	return
}
