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

func NewToolManager(
	opts ...Option,
) (
	toolManager *ToolManager,
	err error,
) {
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

func WithKind(
	toolKind constants.ToolKind,
) (option Option) {
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

func WithName(
	toolKind constants.ToolKind,
) (option Option) {
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

func WithVersion(
	toolKind constants.ToolKind,
) (option Option) {
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
		url                    func(tool string, version string) (toolUrl string)
	)

	url = func(tool string, version string) (toolUrl string) {
		toolUrl = createHashicorpUrl(tool, version, currentOs, currentArch)

		return
	}

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.Url = url(constants.PACKER, constants.PACKER_VERSION)
		case constants.Terraform:
			toolManager.Url = url(constants.TERRAFORM, constants.TERRAFORM_VERSION)
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
	toolKind constants.ToolKind,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithAbsPathToExecutable").
				With("toolKind", toolKind)

		absPathToExecutable func(toolDir string) (apte string)
	)

	absPathToExecutable = func(toolDir string) (apte string) {
		apte = filepath.Join(
			kumoExecAbsPath,
			constants.DEPENDENCIES_DIR_NAME,
			toolDir,
			fmt.Sprintf("%s.exe", toolDir),
		)

		return
	}

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPathTo.Executable = absPathToExecutable(constants.PACKER)

		case constants.Terraform:
			toolManager.AbsPathTo.Executable = absPathToExecutable(constants.TERRAFORM)

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
	toolKind constants.ToolKind,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithAbsPathToDirRun").
				With("toolKind", toolKind).
				With("cloud", cloud).
				With("kumoExecAbsPath", kumoExecAbsPath)

		absPathToDirRun func(toolDir string) (aptdr string)
	)

	absPathToDirRun = func(toolDir string) (aptdr string) {
		aptdr = filepath.Join(
			kumoExecAbsPath,
			toolDir,
			cloud.Name,
		)

		return
	}

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPathTo.Dir.Run = absPathToDirRun(constants.PACKER)

		case constants.Terraform:
			toolManager.AbsPathTo.Dir.Run = absPathToDirRun(constants.TERRAFORM)

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
	toolKind constants.ToolKind,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithAbsPathToDirPlugins").
				With("toolKind", toolKind).
				With("cloud", cloud).
				With("kumoExecAbsPath", kumoExecAbsPath)

		absPathToDirPlugins func(toolDir string) (aptdp string)
	)

	absPathToDirPlugins = func(toolDir string) (aptdp string) {
		aptdp = filepath.Join(
			kumoExecAbsPath,
			toolDir,
			cloud.Name,
			constants.PLUGINS_DIR_NAME,
		)

		return
	}

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPathTo.Dir.Plugins = absPathToDirPlugins(constants.PACKER)

		case constants.Terraform:
			toolManager.AbsPathTo.Dir.Plugins = absPathToDirPlugins(constants.TERRAFORM)

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
	option = func(toolManager *ToolManager) (err error) {
		toolManager.AbsPathTo.Dir.Initial = kumoExecAbsPath

		return
	}

	return
}

func WithAbsPathToTemplateFileMerged(
	toolKind constants.ToolKind,
) (option Option) {
	option = func(toolManager *ToolManager) (err error) {
		toolManager.AbsPathTo.TemplateFile.Merged = filepath.Join(
			constants.TEMPLATES_DIR_NAME,
			constants.MERGED_TEMPLATE,
		)

		return
	}

	return
}

func WithAbsPathToTemplateFileGeneral(
	toolKind constants.ToolKind,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithAbsPathToTemplateFileGeneral").
				With("toolKind", toolKind).
				With("kumoExecAbsPath", kumoExecAbsPath)

		templateGeneralPath func(toolDir string) (tgpath string)
	)

	templateGeneralPath = func(toolDir string) (tgpath string) {
		tgpath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			toolDir,
			constants.GENERAL_DIR_NAME,
			constants.GENERAL_TEMPLATE,
		)

		return
	}

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPathTo.TemplateFile.General = templateGeneralPath(constants.PACKER)

		case constants.Terraform:
			toolManager.AbsPathTo.TemplateFile.General = templateGeneralPath(constants.TERRAFORM)

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
	toolKind constants.ToolKind,
	cloud cloud.Cloud,
	kumoExecAbsPath string,
) (option Option) {

	var (
		oopsBuilder = oops.
				Code("WithAbsPathToTemplateFileCloud").
				With("toolKind", toolKind).
				With("cloud", cloud).
				With("kumoExecAbsPath", kumoExecAbsPath)

		absPathToTemplateCloud func(toolDir string) (tcpath string)
	)

	absPathToTemplateCloud = func(toolDir string) (tcpath string) {
		tcpath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			toolDir,
			fmt.Sprintf("%s.tmpl", cloud.Name),
		)

		return
	}

	option = func(toolManager *ToolManager) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.AbsPathTo.TemplateFile.Cloud = absPathToTemplateCloud(constants.PACKER)

		case constants.Terraform:
			toolManager.AbsPathTo.TemplateFile.Cloud = absPathToTemplateCloud(constants.TERRAFORM)

		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func (tm *ToolManager) SetPluginsPath(
	environmentSetter EnvironmentSetterF,
) (err error) {
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

func (tm *ToolManager) UnsetPluginsPath(
	environmentUnsetter EnvironmentUnsetterF,
) (err error) {
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

func (tm *ToolManager) ChangeToInitialDir(
	dirChanger DirChangerF,
) (err error) {
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

func (tm *ToolManager) ChangeToRunDir(
	dirChanger DirChangerF,
) (err error) {
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
	Executable   string
	Dir          *Dir
	TemplateFile *TemplateFile
}

type Dir struct {
	Plugins string
	Run     string
	Initial string
}

type TemplateFile struct {
	General string
	Cloud   string
	Merged  string
}

type Option func(*ToolManager) error

type EnvironmentSetterF func(key string, value string) error
type EnvironmentUnsetterF func(key string) error
type DirChangerF func(dir string) error
