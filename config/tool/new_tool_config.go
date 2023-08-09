package tool

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/config/cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/utils/host"
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
)

func NewToolConfig(
	opts ...Option,
) (
	toolManager *ToolConfig,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("NewTool").
				With("opts", opts)

		option Option
	)

	toolManager = &ToolConfig{}
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

	option = func(toolManager *ToolConfig) (err error) {
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

	option = func(toolManager *ToolConfig) (err error) {
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

	option = func(toolManager *ToolConfig) (err error) {
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
		url                    = func(tool string, version string) (toolUrl string) {
			toolUrl = createHashicorpUrl(tool, version, currentOs, currentArch)

			return
		}
	)

	option = func(toolManager *ToolConfig) (err error) {
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

	option = func(toolManager *ToolConfig) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.absPath.executable = absPathToExecutable(constants.PACKER)

		case constants.Terraform:
			toolManager.absPath.executable = absPathToExecutable(constants.TERRAFORM)

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
	cloud cloud.CloudConfig,
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

	option = func(toolManager *ToolConfig) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.absPath.dir.run = absPathToDirRun(constants.PACKER)

		case constants.Terraform:
			toolManager.absPath.dir.run = absPathToDirRun(constants.TERRAFORM)

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
	cloud cloud.CloudConfig,
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

	option = func(toolManager *ToolConfig) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.absPath.dir.plugins = absPathToDirPlugins(constants.PACKER)

		case constants.Terraform:
			toolManager.absPath.dir.plugins = absPathToDirPlugins(constants.TERRAFORM)

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
	option = func(toolManager *ToolConfig) (err error) {
		toolManager.absPath.dir.initial = kumoExecAbsPath

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

	option = func(toolManager *ToolConfig) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.absPath.templateFile.general = absPathToTemplateFileGeneral(constants.PACKER)

		case constants.Terraform:
			toolManager.absPath.templateFile.general = absPathToTemplateFileGeneral(constants.TERRAFORM)

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
	cloud cloud.CloudConfig,
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

	option = func(toolManager *ToolConfig) (err error) {
		switch toolKind {
		case constants.Packer:
			toolManager.absPath.templateFile.cloud = absPathToTemplateFileCloud(constants.PACKER)

		case constants.Terraform:
			toolManager.absPath.templateFile.cloud = absPathToTemplateFileCloud(constants.TERRAFORM)

		default:
			err = oopsBuilder.
				Errorf("Unknown tool kind: %d", toolKind)
			return
		}

		return
	}

	return
}

func (tm *ToolConfig) SetPluginsPath(
	os_Setenv func(key string, value string) error,
) (err error) {
	var (
		oopsBuilder = oops.
			Code("SetPluginsDir")
	)

	if err = os_Setenv(constants.PACKER_PLUGIN_PATH, tm.absPath.dir.plugins); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to set plugins dir '%s'", tm.absPath.dir.plugins)
		return
	}

	return
}

func (tm *ToolConfig) UnsetPluginsPath(
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

func (tm *ToolConfig) ChangeToInitialDir(
	os_Chdir DirChangerF,
) (err error) {
	var (
		oopsBuilder = oops.
			Code("ChangeToInitialDir").
			With("dirChanger", os_Chdir)
	)

	if err = os_Chdir(tm.absPath.dir.initial); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to change to initial dir '%s'", tm.absPath.dir.initial)
		return
	}

	return
}

func (tm *ToolConfig) ChangeToRunDir(
	os_Chdir DirChangerF,
) (err error) {
	var (
		oopsBuilder = oops.
			Code("ChangeToRunDir").
			With("dirChanger", os_Chdir)
	)

	if err = os_Chdir(tm.absPath.dir.run); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to change to run dir '%s'", tm.absPath.dir.run)
		return
	}

	return
}

func (tc *ToolConfig) Kind() (kind constants.ToolKind) {
	kind = tc.kind

	return
}

func (tc *ToolConfig) Name() (name string) {
	name = tc.name

	return
}

func (tc *ToolConfig) Version() (version string) {
	version = tc.version

	return
}

func (tc *ToolConfig) Url() (url string) {
	url = tc.url

	return
}

func (tc *ToolConfig) AbsPath() (absPath *AbsPath) {
	absPath = tc.absPath

	return
}

type ToolConfig struct {
	kind    constants.ToolKind
	name    string
	version string
	url     string
	absPath *AbsPath
}

type AbsPath struct {
	executable   string
	dir          *Dir
	templateFile *TemplateFileCombo
}

func (ap *AbsPath) Executable() (executable string) {
	executable = ap.executable

	return
}

func (ap *AbsPath) Dir() (dir *Dir) {
	dir = ap.dir

	return
}

func (ap *AbsPath) TemplateFile() (templateFile *TemplateFileCombo) {
	templateFile = ap.templateFile

	return
}

type Dir struct {
	plugins string
	run     string
	initial string
}

func (d *Dir) Plugins() (plugins string) {
	plugins = d.plugins

	return
}

func (d *Dir) Run() (run string) {
	run = d.run

	return
}

func (d *Dir) Initial() (initial string) {
	initial = d.initial

	return
}

type TemplateFileCombo struct {
	general string
	cloud   string
}

func (tf *TemplateFileCombo) General() (general string) {
	general = tf.general

	return
}

func (tf *TemplateFileCombo) Cloud() (cloud string) {
	cloud = tf.cloud

	return
}

type Option func(*ToolConfig) error

type DirChangerF func(dir string) error
