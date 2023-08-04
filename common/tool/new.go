package tool

import (
	"fmt"
	"os"
	"path/filepath"

	common_cloud_constants "github.com/ed3899/kumo/common/cloud/constants"
	common_cloud_interfaces "github.com/ed3899/kumo/common/cloud/interfaces"
	"github.com/ed3899/kumo/common/dirs"
	common_templates_constants "github.com/ed3899/kumo/common/templates/constants"
	common_templates_interfaces "github.com/ed3899/kumo/common/templates/interfaces"
	"github.com/ed3899/kumo/common/tool/constants"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
)

type Info struct {
	absPath     string
	environment common_templates_interfaces.Environment
}

type Templates struct {
	general *Info
	cloud   *Info
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
	)

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
				general: &Info{
					absPath: filepath.Join(
						kumoExecAbsPath,
						dirs.TEMPLATES_DIR_NAME,
						constants.PACKER_NAME,
						common_templates_constants.PACKER_GENERAL_TEMPLATE_NAME,
					),
				},
			},
		}

	case constants.Terraform:
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
				general: &Info{
					absPath: filepath.Join(
						kumoExecAbsPath,
						dirs.TEMPLATES_DIR_NAME,
						constants.TERRAFORM_NAME,
						common_templates_constants.TERRAFORM_GENERAL_TEMPLATE_NAME,
					),
				},
			},
		}

	default:
		err = oopsBuilder.
			Errorf("Tool '%v' not supported", toolKind)
		return
	}

	switch cloud.Kind() {
	case common_cloud_constants.AWS:
	default:
		err = oopsBuilder.
			Wrapf(err, "Cloud '%v' not supported", cloud.Kind())
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
