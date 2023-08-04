package tool

import (
	"fmt"
	"os"
	"path/filepath"

	cloud_constants "github.com/ed3899/kumo/common/cloud/constants"
	common_cloud_interfaces "github.com/ed3899/kumo/common/cloud/interfaces"
	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/tool/constants"
	packer_aws "github.com/ed3899/kumo/common/tool/environments/packer/aws"
	packer_general "github.com/ed3899/kumo/common/tool/environments/packer/general"
	terraform_aws "github.com/ed3899/kumo/common/tool/environments/terraform/aws"
	terraform_general "github.com/ed3899/kumo/common/tool/environments/terraform/general"
	"github.com/ed3899/kumo/common/tool/packerManifest"
	"github.com/ed3899/kumo/common/tool/interfaces"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Data struct {
	absPath     string
	environment interfaces.Environment
}

type Templates struct {
	general *Data
	cloud   *Data
}

type Tool struct {
	kind              constants.Kind
	name              string
	version           string
	executableAbsPath string
	runDir            string
	pluginDir         string
	templates         *Templates
}

func New(toolKind constants.Kind, cloud common_cloud_interfaces.Cloud, kumoExecAbsPath string) (toolConfig *Tool, err error) {
	var (
		oopsBuilder = oops.
				Code("new_tool_setup_failed").
				With("tool", toolKind)
		logger, _ = zap.NewProduction()
		
		publicIp string
		pickedIp string
	)

	defer logger.Sync()

	switch toolKind {
	case constants.Packer:
		toolConfig = &Tool{
			kind:    constants.Packer,
			name:    constants.PACKER_NAME,
			version: constants.PACKER_VERSION,
			executableAbsPath: filepath.Join(
				kumoExecAbsPath,
				dirs.DEPENDENCIES_DIR_NAME,
				constants.PACKER_NAME,
				fmt.Sprintf("%s.exe", constants.PACKER_NAME),
			),
			runDir: filepath.Join(
				kumoExecAbsPath,
				constants.PACKER_NAME,
				cloud.Name(),
			),
			pluginDir: filepath.Join(
				kumoExecAbsPath,
				constants.PACKER_NAME,
				cloud.Name(),
				dirs.PLUGINS_DIR_NAME,
			),
			templates: &Templates{
				general: &Data{
					absPath: filepath.Join(
						kumoExecAbsPath,
						dirs.TEMPLATES_DIR_NAME,
						constants.PACKER_NAME,
						constants.PACKER_GENERAL_TEMPLATE_NAME,
					),
					environment: &packer_general.Environment{
						Required: &packer_general.Required{
							GIT_USERNAME: viper.GetString("Git.Username"),
							GIT_EMAIL:    viper.GetString("Git.Email"),
							ANSIBLE_TAGS: viper.GetStringSlice("AMI.Tools"),
						},
						Optional: &packer_general.Optional{
							GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC: viper.GetString("GitHub.PersonalAccessTokenClassic"),
						},
					},
				},
			},
		}

		switch cloud.Kind() {
		case cloud_constants.AWS:
			toolConfig.templates.cloud = &Data{
				absPath: filepath.Join(
					kumoExecAbsPath,
					dirs.TEMPLATES_DIR_NAME,
					constants.PACKER_NAME,
					cloud_constants.AWS_NAME,
					constants.PACKER_AWS_TEMPLATE_NAME,
				),
				environment: &packer_aws.Environment{
					Required: &packer_aws.Required{
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
					},
				},
			}

		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s not supported", cloud.Name())
			return
		}

	case constants.Terraform:
		if publicIp, err = utils.GetPublicIp(); err != nil {
			logger.Warn(
				"Failed to get public IP, using default ip",
				zap.String("error", err.Error()),
				zap.String("defaultIp", constants.TERRAFORM_DEFAULT_ALLOWED_IP),
			)
			pickedIp = constants.TERRAFORM_DEFAULT_ALLOWED_IP
		} else {
			pickedIp = publicIp
		}

		toolConfig = &Tool{
			kind:    constants.Terraform,
			name:    constants.TERRAFORM_NAME,
			version: constants.TERRAFORM_VERSION,
			executableAbsPath: filepath.Join(
				kumoExecAbsPath,
				dirs.DEPENDENCIES_DIR_NAME,
				constants.TERRAFORM_NAME,
				fmt.Sprintf("%s.exe", constants.TERRAFORM_NAME),
			),
			runDir: filepath.Join(
				kumoExecAbsPath,
				constants.TERRAFORM_NAME,
				cloud.Name(),
			),
			pluginDir: filepath.Join(
				kumoExecAbsPath,
				constants.TERRAFORM_NAME,
				cloud.Name(),
				dirs.PLUGINS_DIR_NAME,
			),
			templates: &Templates{
				general: &Data{
					absPath: filepath.Join(
						kumoExecAbsPath,
						dirs.TEMPLATES_DIR_NAME,
						constants.TERRAFORM_NAME,
						constants.TERRAFORM_GENERAL_TEMPLATE_NAME,
					),
					environment: &terraform_general.Environment{
						Required: &terraform_general.Required{
							ALLOWED_IP: pickedIp,
						},
					},
				},
			},
		}

		switch cloud.Kind() {
		case cloud_constants.AWS:
			toolConfig.templates.cloud = &Data{
				absPath: filepath.Join(
					kumoExecAbsPath,
					dirs.TEMPLATES_DIR_NAME,
					constants.TERRAFORM_NAME,
					cloud_constants.AWS_NAME,
					constants.TERRAFORM_AWS_TEMPLATE_NAME,
				),
				environment: &terraform_aws.Environment{
					Required: &terraform_aws.Required{
						AWS_REGION:        viper.GetString("AWS.Region"),
						AWS_INSTANCE_TYPE: viper.GetString("AWS.EC2.Instance.Type"),
						AMI_ID:            pickedAmiId,
						KEY_NAME:          ssh.KEY_NAME,
						SSH_PORT:          ssh.SSH_PORT,
						IP_FILE_NAME:      ssh.IP_FILE_NAME,
						USERNAME:          viper.GetString("AMI.User"),
					},
					Optional: &terraform_aws.Optional{
						AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
						AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
					},
				},
			}
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s not supported", cloud.Name())
			return
		}

	default:
		err = oopsBuilder.
			Errorf("Tool '%v' not supported", toolKind)
		return
	}

	return
}

func (t *Tool) Kind() (toolKind constants.Kind) {
	return t.kind
}

func (t *Tool) Name() (toolName string) {
	return t.name
}

func (t *Tool) SetPluginPath() (err error) {
	var (
		oopsBuilder = oops.
			Code("set_plugin_path_failed")
	)

	if err = os.Setenv(constants.PACKER_PLUGIN_PATH_NAME, t.pluginDir); err != nil {
		err = oopsBuilder.
			With("pluginDir", t.pluginDir).
			Wrapf(err, "Error occurred while setting plugin path for %s", t.name)
		return
	}

	return
}

func (t *Tool) UnsetPluginPath() (err error) {
	var (
		oopsBuilder = oops.
			Code("unset_plugin_path_failed")
	)

	if err = os.Unsetenv(constants.PACKER_PLUGIN_PATH_NAME); err != nil {
		err = oopsBuilder.
			With("pluginDir", t.pluginDir).
			Wrapf(err, "Error occurred while unsetting plugin path for %s", t.name)
		return
	}

	return
}

func (t *Tool) ExecutableName() (toolExecutableName string) {
	return fmt.Sprintf("%s.exe", t.name)
}

func (t *Tool) Version() (toolVersion string) {
	return t.version
}

func (t *Tool) RunDir() (toolDir string) {
	return t.runDir
}

func (t *Tool) Url() (toolUrl string) {
	var (
		currentOs, currentArch = utils.GetCurrentHostSpecs()
	)

	return utils.CreateHashicorpURL(t.name, t.version, currentOs, currentArch)
}
