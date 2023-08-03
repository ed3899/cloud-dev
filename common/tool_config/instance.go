package tool_config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/common/cloud_config"
	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
)

type Tool struct {
	kind               Kind
	name                string
	version             string
	dir                 string
	dependenciesDirName string
	initialDir          string
}

func New(toolType Kind, cloud cloud_config.CloudI) (toolConfig *Tool, err error) {
	var (
		oopsBuilder = oops.
				Code("new_tool_setup_failed").
				With("tool", toolType)

		cwd string
	)

	if cwd, err = os.Getwd(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while getting current working directory")
		return
	}

	switch toolType {
	case Packer:
		toolConfig = &Tool{
			dependenciesDirName: dirs.DEPENDENCIES_DIR_NAME,
			kind:               Packer,
			name:                PACKER_NAME,
			version:             PACKER_VERSION,
			initialDir:          cwd,
			dir:                 filepath.Join(PACKER_NAME, cloud.Name()),
		}

	case Terraform:
		toolConfig = &Tool{
			dependenciesDirName: dirs.DEPENDENCIES_DIR_NAME,
			kind:               Terraform,
			name:                TERRAFORM_NAME,
			version:             TERRAFORM_VERSION,
			initialDir:          cwd,
			dir:                 filepath.Join(TERRAFORM_NAME, cloud.Name()),
		}

	default:
		err = oopsBuilder.
			Errorf("Tool '%v' not supported", toolType)
		return
	}

	return
}

func (t *Tool) DependenciesDirName() (dependenciesDirName string) {
	return t.dependenciesDirName
}

func (t *Tool) Type() (toolType Kind) {
	return t.kind
}

func (t *Tool) Name() (toolName string) {
	return t.name
}

func (t *Tool) ZipName() (toolZipName string) {
	return fmt.Sprintf("%s.zip", t.name)
}

func (t *Tool) ZipAbsPath() (toolZipAbsPath string) {
	return filepath.Join(t.dependenciesDirName, t.name, fmt.Sprintf("%s.zip", t.name))
}

func (t *Tool) ZipContentLength() (toolZipContentLength int64, err error) {
	var (
		oopsBuilder = oops.
				Code("get_zip_content_length_failed")
		currentOs, currentArch = utils.GetCurrentHostSpecs()
		url                    = utils.CreateHashicorpURL(t.name, t.version, currentOs, currentArch)
	)

	if toolZipContentLength, err = utils.GetContentLength(url); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get content length for: %s", url)
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

func (t *Tool) InitialDir() (initialDir string) {
	return t.initialDir
}

func (t *Tool) Dir() (toolDir string) {
	return t.dir
}

func (t *Tool) GetUrl() (toolUrl string) {
	var (
		currentOs, currentArch = utils.GetCurrentHostSpecs()
	)

	return utils.CreateHashicorpURL(t.name, t.version, currentOs, currentArch)
}

func (t *Tool) GoInitialDir() (err error) {
	var (
		oopsBuilder = oops.
			Code("go_initial_dir_failed")
	)

	if err = os.Chdir(t.initialDir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while changing directory to %s", t.initialDir)
		return
	}

	return
}

func (t *Tool) GoDir() (err error) {
	var (
		oopsBuilder = oops.
			Code("go_target_dir_failed")
	)

	if err = os.Chdir(t.dir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while changing directory to %s", t.dir)
		return
	}

	return
}
