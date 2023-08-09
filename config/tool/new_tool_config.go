package tool

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/common/interfaces"
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

func WithKind(
	toolKind constants.ToolKind,
) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithKind").
			With("toolKind", toolKind)
	)

	option = func(toolManager *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.kind = constants.Packer
		case constants.Terraform:
			toolManager.kind = constants.Terraform
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
	toolKind constants.ToolKind,
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
	toolKind constants.ToolKind,
	createHashicorpUrl url.CreateHashicorpURLF,
	getCurrentHostSpecs host.GetCurrentHostSpecsF,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithUrl").
				With("toolKind", toolKind)
		currentOs, currentArch = getCurrentHostSpecs()
		url                    = func(tool string, version string) (toolUrl ToolUrl) {
			toolUrl = ToolUrl(createHashicorpUrl(tool, version, currentOs, currentArch))

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
	toolKind constants.ToolKind,
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
	toolKind constants.ToolKind,
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
	toolKind constants.ToolKind,
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
	toolKind constants.ToolKind,
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

func (tm *Tool) SetPluginsPath(
	os_Setenv func(key string, value string) error,
) (err error) {
	var (
		oopsBuilder = oops.
			Code("SetPluginsDir")
	)

	if err = os_Setenv(constants.PACKER_PLUGIN_PATH, tm.AbsPath.Dir.Plugins); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to set plugins dir '%s'", tm.AbsPath.Dir.Plugins)
		return
	}

	return
}

func (tm *Tool) UnsetPluginsPath(
	os_Unset func(key string) error,
) (err error) {
	var (
		oopsBuilder = oops.
			Code("UnsetPluginsDir")
	)

	if err = os_Unset(constants.PACKER_PLUGIN_PATH); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to unset plugins dir '%s'", constants.PACKER_PLUGIN_PATH)
		return
	}

	return
}

func (tm *Tool) ChangeToInitialDir(
	os_Chdir DirChangerF,
) (err error) {
	var (
		oopsBuilder = oops.
			Code("ChangeToInitialDir").
			With("dirChanger", os_Chdir)
	)

	if err = os_Chdir(tm.AbsPath.Dir.Initial); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to change to initial dir '%s'", tm.AbsPath.Dir.Initial)
		return
	}

	return
}

func (tm *Tool) ChangeToRunDir(
	os_Chdir DirChangerF,
) (err error) {
	var (
		oopsBuilder = oops.
			Code("ChangeToRunDir").
			With("dirChanger", os_Chdir)
	)

	if err = os_Chdir(tm.AbsPath.Dir.Run); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to change to run dir '%s'", tm.AbsPath.Dir.Run)
		return
	}

	return
}

func (tc *Tool) Kind() (toolKind constants.ToolKind) {
	toolKind = tc.kind
	return
}

func (tc *Tool) Name() (toolName ToolName) {
	toolName = tc.name
	return
}

func (tc *Tool) Version() (toolVersion ToolVersion) {
	toolVersion = tc.version
	return
}

func (tc *Tool) Url() (toolUrl ToolUrl) {
	toolUrl = tc.url
	return
}

type DirChangerF func(dir string) error

type Option func(*Tool) error

type ToolName string

func (tn ToolName) String() (toolName string) {
	toolName = string(tn)
	return
}

type ToolVersion string

func (tv ToolVersion) String() (toolVersion string) {
	toolVersion = string(tv)
	return
}

type ToolUrl string

func (tu ToolUrl) String() (toolUrl string) {
	toolUrl = string(tu)
	return
}

type Tool struct {
	kind    constants.ToolKind
	name    ToolName
	version ToolVersion
	url     ToolUrl
	AbsPath *ToolAbsPath
}

type ToolAbsPath struct {
	Executable   string
	Dir          *ToolDir
	TemplateFile *TemplateFileCombo
}

type TemplateFileCombo struct {
	General string
	Cloud   string
}

type ToolI interface {
	interfaces.KindGetter[constants.ToolKind]
	interfaces.NameGetter[ToolName]
	interfaces.VersionGetter[ToolVersion]
	interfaces.UrlGetter[ToolUrl]
}
