package tool_manager

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/utils/host"
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
)

func NewToolManager(opts ...Option) (toolManager *ToolManager, err error) {
	var (
		oopsBuilder = oops.
				Code("NewTool").
				With("opts", opts)

		option Option
	)

	toolManager = &ToolManager{}
	for _, option = range opts {
		if err = option(toolManager); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Option %v", option)
			return
		}
	}

	return
}

func WithKind(toolKind constants.ToolKind) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithKind").
			With("toolKind", toolKind)
	)

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.Kind = constants.Packer
		case constants.Terraform:
			toolManager.Kind = constants.Terraform
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithName(toolKind constants.ToolKind) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithName").
			With("toolKind", toolKind)
	)

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.Name = constants.PACKER
		case constants.Terraform:
			toolManager.Name = constants.TERRAFORM
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithVersion(toolKind constants.ToolKind) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithVersion").
			With("toolKind", toolKind)
	)

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.Version = constants.PACKER_VERSION
		case constants.Terraform:
			toolManager.Version = constants.TERRAFORM_VERSION
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
	toolKind constants.ToolKind,
	createHashicorpUrl url.CreateHashicorpURLF,
	getCurrentHostSpecs host.GetCurrentHostSpecsF,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithUrl").
				With("toolKind", toolKind)
		currentOs, currentArch = getCurrentHostSpecs()
	)

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.Url = createHashicorpUrl(constants.PACKER, constants.PACKER_VERSION, currentOs, currentArch)
		case constants.Terraform:
			toolManager.Url = createHashicorpUrl(constants.TERRAFORM, constants.TERRAFORM_VERSION, currentOs, currentArch)
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithExecutableAbsPath(toolKind constants.ToolKind, kumoExecAbsPath string) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithExecutableAbsPath").
			With("toolKind", toolKind)
	)

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPathTo.Executable = filepath.Join(
				kumoExecAbsPath,
				constants.DEPENDENCIES_DIR_NAME,
				constants.PACKER,
				fmt.Sprintf("%s.exe", constants.PACKER),
			)
		case constants.Terraform:
			toolManager.AbsPathTo.Executable = filepath.Join(
				kumoExecAbsPath,
				constants.DEPENDENCIES_DIR_NAME,
				constants.TERRAFORM,
				fmt.Sprintf("%s.exe", constants.TERRAFORM),
			)
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithRunDirAbsPath(cloud cloud.Cloud, toolKind constants.ToolKind, kumoExecAbsPath string) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithRunDir").
			With("toolKind", toolKind).
			With("cloud", cloud).
			With("kumoExecAbsPath", kumoExecAbsPath)
	)

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPathTo.Dir.Run = filepath.Join(
				kumoExecAbsPath,
				constants.PACKER,
				cloud.Name,
			)
		case constants.Terraform:
			toolManager.AbsPathTo.Dir.Run = filepath.Join(
				kumoExecAbsPath,
				constants.TERRAFORM,
				cloud.Name,
			)
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithPluginsDir(cloud cloud.Cloud, toolKind constants.ToolKind, kumoExecAbsPath string) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithPluginsDir").
			With("toolKind", toolKind).
			With("cloud", cloud).
			With("kumoExecAbsPath", kumoExecAbsPath)
	)

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPathTo.Dir.Plugins = filepath.Join(
				kumoExecAbsPath,
				constants.PACKER,
				cloud.Name,
				constants.PLUGINS_DIR_NAME,
			)
		case constants.Terraform:
			toolManager.AbsPathTo.Dir.Plugins = filepath.Join(
				kumoExecAbsPath,
				constants.TERRAFORM,
				cloud.Name,
				constants.PLUGINS_DIR_NAME,
			)
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithInitialDirAbsPath(kumoExecAbsPath string) (option Option) {
	option = func(toolManager *ToolManager) (err error) {
		toolManager.AbsPathTo.Dir.Initial = kumoExecAbsPath

		return
	}

	return
}

func WithTempMergedTemplateFileName(toolKind constants.ToolKind) (option Option) {
	var (
		oopsBuilder = oops.
			Code("TempMergedTemplateFileName").
			With("toolKind", toolKind)
	)

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPathTo.Template.Merged = filepath.Join(constants.TEMPLATES_DIR_NAME, constants.MERGED_TEMPLATE)
		case constants.Terraform:
			toolManager.AbsPathTo.Template.Merged = filepath.Join(constants.TEMPLATES_DIR_NAME, constants.MERGED_TEMPLATE)
		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func WithAbsPathToGeneralTemplate(
	toolKind constants.ToolKind,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithAbsPathToGeneralTemplateFor").
			With("toolKind", toolKind).
			With("kumoExecAbsPath", kumoExecAbsPath)
	)

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPathTo.Template.General = filepath.Join(
				kumoExecAbsPath,
				constants.TEMPLATES_DIR_NAME,
				constants.PACKER,
				constants.GENERAL_DIR_NAME,
				constants.GENERAL_TEMPLATE,
			)

		case constants.Terraform:
			toolManager.AbsPathTo.Template.General = filepath.Join(
				kumoExecAbsPath,
				constants.TEMPLATES_DIR_NAME,
				constants.TERRAFORM,
				constants.GENERAL_DIR_NAME,
				constants.GENERAL_TEMPLATE,
			)

		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}
	return
}

func WithAbsPathToCloudTemplate(
	toolKind constants.ToolKind,
	cloud cloud.Cloud,
	kumoExecAbsPath string,
) (option Option) {

	var (
		oopsBuilder = oops.
			Code("WithAbsPathToCloudTemplateFor").
			With("toolKind", toolKind).
			With("cloud", cloud).
			With("kumoExecAbsPath", kumoExecAbsPath)
	)

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPathTo.Template.Cloud = filepath.Join(
				kumoExecAbsPath,
				constants.TEMPLATES_DIR_NAME,
				constants.PACKER,
				fmt.Sprintf("%s.tmpl", cloud.Name),
			)

		case constants.Terraform:
			toolManager.AbsPathTo.Template.Cloud = filepath.Join(
				kumoExecAbsPath,
				constants.TEMPLATES_DIR_NAME,
				constants.TERRAFORM,
				fmt.Sprintf("%s.tmpl", cloud.Name),
			)

		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func (tm *ToolManager) SetPluginsPathWith(environmentSetter EnvironmentSetterF) (err error) {
	var (
		oopsBuilder = oops.
			Code("SetPluginsDir")
	)

	if err = environmentSetter(constants.PACKER_PLUGIN_PATH, tm.AbsPathTo.Dir.Plugins); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to set plugins dir '%s'", tm.AbsPathTo.Dir.Plugins)
		return
	}

	return
}

func (tm *ToolManager) UnsetPluginsPathWith(environmentUnsetter EnvironmentUnsetterF) (err error) {
	var (
		oopsBuilder = oops.
			Code("UnsetPluginsDir")
	)

	if err = environmentUnsetter(constants.PACKER_PLUGIN_PATH); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to unset plugins dir '%s'", constants.PACKER_PLUGIN_PATH)
		return
	}

	return
}

func (tm *ToolManager) ChangeToInitialDirWith(dirChanger DirChangerF) (err error) {
	var (
		oopsBuilder = oops.
			Code("ChangeToInitialDir").
			With("dirChanger", dirChanger)
	)

	if err = dirChanger(tm.AbsPathTo.Dir.Initial); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to change to initial dir '%s'", tm.AbsPathTo.Dir.Initial)
		return
	}

	return
}

func (tm *ToolManager) ChangeToRunDirWith(dirChanger DirChangerF) (err error) {
	var (
		oopsBuilder = oops.
			Code("ChangeToRunDir").
			With("dirChanger", dirChanger)
	)

	if err = dirChanger(tm.AbsPathTo.Dir.Run); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to change to run dir '%s'", tm.AbsPathTo.Dir.Run)
		return
	}

	return
}

type ToolManager struct {
	Kind      constants.ToolKind
	Name      string
	Version   string
	Url       string
	AbsPathTo *AbsPathTo
}

type AbsPathTo struct {
	Executable string
	Dir        *Dir
	Template   *Template
}

type Dir struct {
	Plugins string
	Run     string
	Initial string
}

type Template struct {
	General string
	Cloud   string
	Merged  string
}

type Option func(*ToolManager) error

type EnvironmentSetterF func(key string, value string) error
type EnvironmentUnsetterF func(key string) error
type DirChangerF func(dir string) error
