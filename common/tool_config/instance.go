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
	kind              Kind
	name              string
	version           string
	executableAbsPath string
	runDir            string
}

func New(toolKind Kind, cloud cloud_config.CloudI, kumoExecAbsPath string) (toolConfig *Tool, err error) {
	var (
		oopsBuilder = oops.
			Code("new_tool_setup_failed").
			With("tool", toolKind)
	)

	switch toolKind {
	case Packer:
		toolConfig = &Tool{
			kind:    Packer,
			name:    PACKER_NAME,
			version: PACKER_VERSION,
			executableAbsPath: filepath.Join(
				kumoExecAbsPath,
				dirs.DEPENDENCIES_DIR_NAME,
				PACKER_NAME,
				fmt.Sprintf("%s.exe", PACKER_NAME),
			),
			runDir: filepath.Join(
				kumoExecAbsPath,
				PACKER_NAME,
				cloud.Name(),
			),
		}

	case Terraform:
		toolConfig = &Tool{
			kind:    Terraform,
			name:    TERRAFORM_NAME,
			version: TERRAFORM_VERSION,
			executableAbsPath: filepath.Join(
				kumoExecAbsPath,
				dirs.DEPENDENCIES_DIR_NAME,
				TERRAFORM_NAME,
				fmt.Sprintf("%s.exe", TERRAFORM_NAME),
			),
			runDir: filepath.Join(
				kumoExecAbsPath,
				TERRAFORM_NAME,
				cloud.Name(),
			),
		}

	default:
		err = oopsBuilder.
			Errorf("Tool '%v' not supported", toolKind)
		return
	}

	return
}

func (t *Tool) Kind() (toolKind Kind) {
	return t.kind
}

func (t *Tool) Name() (toolName string) {
	return t.name
}

// func (t *Tool) ZipName() (toolZipName string) {
// 	return fmt.Sprintf("%s.zip", t.name)
// }

// func (t *Tool) ZipAbsPath() (toolZipAbsPath string) {
// 	return filepath.Join(t.dependenciesDirName, t.name, fmt.Sprintf("%s.zip", t.name))
// }

// func (t *Tool) ZipContentLength() (toolZipContentLength int64, err error) {
// 	var (
// 		oopsBuilder = oops.
// 				Code("get_zip_content_length_failed")
// 		currentOs, currentArch = utils.GetCurrentHostSpecs()
// 		url                    = utils.CreateHashicorpURL(t.name, t.version, currentOs, currentArch)
// 	)

// 	if toolZipContentLength, err = utils.GetContentLength(url); err != nil {
// 		err = oopsBuilder.
// 			Wrapf(err, "failed to get content length for: %s", url)
// 		return
// 	}

// 	return
// }

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

