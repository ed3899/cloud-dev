package tool

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/common/alias"
	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/config/cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/utils/host"
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
)

func NewTool(
	opts ...Option,
) (
	tool *Tool,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("NewTool").
				With("opts", opts)

		option Option
	)

	tool = &Tool{}
	for _, option = range opts {
		if err = option(tool); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Option %v", option)
			return
		}
	}

	return
}

type Tool[
	ToolName ~func() string,
	ToolVersion ~func() string,
	ToolUrl ~func() string,
] struct {
	Name    ToolName
	Version ToolVersion
	Url     ToolUrl
}

func WithName[
	ToolNameWithMaybe ~func(iota.Tool) (alias.ToolName, error),
](
	toolKind constants.Tool,
) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithName").
			With("toolKind", toolKind)
	)

	option = func(toolManager *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.name = constants.PACKER
		case constants.Terraform:
			toolManager.name = constants.TERRAFORM
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithVersion(
	toolKind constants.Tool,
) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithVersion").
			With("toolKind", toolKind)
	)

	option = func(toolManager *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.version = constants.PACKER_VERSION
		case constants.Terraform:
			toolManager.version = constants.TERRAFORM_VERSION
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithUrl(
	toolKind constants.Tool,
	createHashicorpUrl url.CreateHashicorpURLF,
	getCurrentHostSpecs host.GetCurrentHostSpecsF,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithUrl").
				With("toolKind", toolKind)
		currentOs, currentArch = getCurrentHostSpecs()
		url                    = func(tool string, version string) (toolUrl ToolUrl) {
			toolUrl = alias.ToolUrl(createHashicorpUrl(tool, version, currentOs, currentArch))

			return
		}
	)

	option = func(toolManager *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.url = url(constants.PACKER, constants.PACKER_VERSION)
		case constants.Terraform:
			toolManager.url = url(constants.TERRAFORM, constants.TERRAFORM_VERSION)
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithAbsPathToExecutable(
	toolKind constants.Tool,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithAbsPathToExecutable").
				With("toolKind", toolKind)

		absPathToExecutable = func(toolDir string) (apte string) {
			apte = filepath.Join(
				kumoExecAbsPath,
				constants.DEPENDENCIES_DIR_NAME,
				toolDir,
				fmt.Sprintf("%s.exe", toolDir),
			)

			return
		}
	)

	option = func(toolManager *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPath.Executable = absPathToExecutable(constants.PACKER)

		case constants.Terraform:
			toolManager.AbsPath.Executable = absPathToExecutable(constants.TERRAFORM)

		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithAbsPathToDirRun(
	cloud cloud.Cloud,
	toolKind constants.Tool,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithAbsPathToDirRun").
				With("toolKind", toolKind).
				With("cloud", cloud).
				With("kumoExecAbsPath", kumoExecAbsPath)

		absPathToDirRun = func(toolDir string) (aptdr string) {
			aptdr = filepath.Join(
				kumoExecAbsPath,
				toolDir,
				cloud.Name,
			)

			return
		}
	)

	option = func(toolManager *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPath.Dir.Run = absPathToDirRun(constants.PACKER)

		case constants.Terraform:
			toolManager.AbsPath.Dir.Run = absPathToDirRun(constants.TERRAFORM)

		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithAbsPathToDirPlugins(
	cloud cloud.Cloud,
	toolKind constants.Tool,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithAbsPathToDirPlugins").
				With("toolKind", toolKind).
				With("cloud", cloud).
				With("kumoExecAbsPath", kumoExecAbsPath)

		absPathToDirPlugins = func(toolDir string) (aptdp string) {
			aptdp = filepath.Join(
				kumoExecAbsPath,
				toolDir,
				cloud.Name,
				constants.PLUGINS_DIR_NAME,
			)

			return
		}
	)

	option = func(toolManager *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPath.Dir.Plugins = absPathToDirPlugins(constants.PACKER)

		case constants.Terraform:
			toolManager.AbsPath.Dir.Plugins = absPathToDirPlugins(constants.TERRAFORM)

		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithAbsPathToDirInitial(
	kumoExecAbsPath string,
) (option Option) {
	option = func(toolManager *Tool) (err error) {
		toolManager.AbsPath.Dir.Initial = kumoExecAbsPath

		return
	}

	return
}

func WithAbsPathToTemplateFileGeneral(
	toolKind constants.Tool,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithAbsPathToTemplateFileGeneral").
				With("toolKind", toolKind).
				With("kumoExecAbsPath", kumoExecAbsPath)

		absPathToTemplateFileGeneral = func(toolDir string) (tgpath string) {
			tgpath = filepath.Join(
				kumoExecAbsPath,
				constants.TEMPLATES_DIR_NAME,
				toolDir,
				constants.GENERAL_DIR_NAME,
				constants.GENERAL_TEMPLATE,
			)

			return
		}
	)

	option = func(toolManager *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPath.TemplateFile.General = absPathToTemplateFileGeneral(constants.PACKER)

		case constants.Terraform:
			toolManager.AbsPath.TemplateFile.General = absPathToTemplateFileGeneral(constants.TERRAFORM)

		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}
	return
}

func WithAbsPathToTemplateFileCloud(
	toolKind constants.Tool,
	cloud cloud.Cloud,
	kumoExecAbsPath string,
) (option Option) {

	var (
		oopsBuilder = oops.
				Code("WithAbsPathToTemplateFileCloud").
				With("toolKind", toolKind).
				With("cloud", cloud).
				With("kumoExecAbsPath", kumoExecAbsPath)

		absPathToTemplateFileCloud = func(toolDir string) (tcpath string) {
			tcpath = filepath.Join(
				kumoExecAbsPath,
				constants.TEMPLATES_DIR_NAME,
				toolDir,
				fmt.Sprintf("%s.tmpl", cloud.Name),
			)

			return
		}
	)

	option = func(toolManager *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPath.TemplateFile.Cloud = absPathToTemplateFileCloud(constants.PACKER)

		case constants.Terraform:
			toolManager.AbsPath.TemplateFile.Cloud = absPathToTemplateFileCloud(constants.TERRAFORM)

		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

type ToolAbsPath struct {
	Executable   string
	Dir          *ToolDir
	TemplateFile *TemplateFileCombo
}

type DirChangerF func(dir string) error

type Option func(*Tool) error
