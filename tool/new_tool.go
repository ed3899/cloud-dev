package tool

import (
	"fmt"
	"path/filepath"

	"github.com/ed3899/kumo/cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/utils/host"
	"github.com/ed3899/kumo/utils/url"
	"github.com/samber/oops"
)

func NewTool(opts ...Option) (tool *Tool, err error) {
	var (
		oopsBuilder = oops.
				Code("NewTool").
				With("opts", opts)

		o Option
	)

	tool = &Tool{}
	for _, o = range opts {
		if err = o(tool); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Option %v", o)
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

	option = func(tool *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			tool.Kind = constants.Packer
		case constants.Terraform:
			tool.Kind = constants.Terraform
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

	option = func(tool *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			tool.Name = constants.PACKER
		case constants.Terraform:
			tool.Name = constants.TERRAFORM
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

	option = func(tool *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			tool.Version = constants.PACKER_VERSION
		case constants.Terraform:
			tool.Version = constants.TERRAFORM_VERSION
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

	option = func(tool *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			tool.Url = createHashicorpUrl(constants.PACKER, constants.PACKER_VERSION, currentOs, currentArch)
		case constants.Terraform:
			tool.Url = createHashicorpUrl(constants.TERRAFORM, constants.TERRAFORM_VERSION, currentOs, currentArch)
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

	option = func(tool *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			tool.ExecutableAbsPath = filepath.Join(
				kumoExecAbsPath,
				constants.DEPENDENCIES_DIR_NAME,
				constants.PACKER,
				fmt.Sprintf("%s.exe", constants.PACKER),
			)
		case constants.Terraform:
			tool.ExecutableAbsPath = filepath.Join(
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

	option = func(tool *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			tool.RunDirAbsPath = filepath.Join(
				kumoExecAbsPath,
				constants.PACKER,
				cloud.Name,
			)
		case constants.Terraform:
			tool.RunDirAbsPath = filepath.Join(
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

	option = func(tool *Tool) (err error) {
		switch toolKind {
		case constants.Packer:
			tool.PluginsDirAbsPath = filepath.Join(
				kumoExecAbsPath,
				constants.PACKER,
				cloud.Name,
				constants.PLUGINS_DIR_NAME,
			)
		case constants.Terraform:
			tool.PluginsDirAbsPath = filepath.Join(
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

type Tool struct {
	Kind              constants.ToolKind
	Name              string
	Version           string
	Url               string
	ExecutableAbsPath string
	RunDirAbsPath     string
	PluginsDirAbsPath string
}

type Option func(*Tool) error
