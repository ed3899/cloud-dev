package packer

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/cloud"
	"github.com/ed3899/kumo/constants"
	environment_packer "github.com/ed3899/kumo/tool_manager/environment/packer"
	template_packer "github.com/ed3899/kumo/tool_manager/template/packer"
	"github.com/ed3899/kumo/utils/file"
	"github.com/samber/oops"
)

type PackerVars struct {
	AbsPath     string
	Environment *environment_packer.PackerEnvironment
	Template    *template_packer.MergedTemplate
}

func NewPackerVars(
	options ...Option,
) (
	packerVars *PackerVars,
) {
	var (
		option Option
	)

	packerVars = &PackerVars{}

	for _, option = range options {
		option(packerVars)
	}

	return
}

func WithEnvironment(
	packerEnvironment *environment_packer.PackerEnvironment,
) (option Option) {
	option = func(packerVars *PackerVars) {
		packerVars.Environment = packerEnvironment

		return
	}

	return
}

func WithTemplate(
	packerTemplateCombo *template_packer.MergedTemplate,
) (option Option) {
	option = func(packerVars *PackerVars) {
		packerVars.Template = packerTemplateCombo

		return
	}

	return
}

func WithAbsPathFor(
	kumoExecAbsPath string,
	cloud cloud.Cloud,
) (option Option) {
	option = func(packerVars *PackerVars) {
		packerVars.AbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.PACKER,
			cloud.Name,
			constants.PACKER_VARS,
		)

		return
	}

	return
}

func (pv *PackerVars) OutputVars(
	kumoExecAbsPath string,
	cloud cloud.Cloud,
	fileMerger file.MergeFilesToF,
	remover func(string) error,
	fileOpener func(string) (*os.File, error),
	fileCreator func(string) (*os.File, error),
	templateParser func(filenames ...string) (*template.Template, error),
) (err error) {
	var (
		oopsBuilder = oops.
				Code("OutputVars").
				With("fileMerger", fileMerger).
				With("remover", remover)
		varsAbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.PACKER,
			cloud.Name,
			constants.PACKER_VARS,
		)
		tempDir = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
		)

		packerVarsFile    *os.File
		mergedFileAbsPath string
		mergedTemplate    *template.Template
	)

	if mergedFileAbsPath, err = fileMerger(tempDir, pv.Template.GeneralTemplateFileAbsPath, pv.Template.CloudTemplateFileAbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to merge general and cloud template '%s' and '%s'", pv.Template.GeneralTemplateFileAbsPath, pv.Template.CloudTemplateFileAbsPath)
		return
	}

	if mergedTemplate, err = templateParser(mergedFileAbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to parse merged template '%s'", mergedFileAbsPath)
		return
	}

	if packerVarsFile, err = fileCreator(varsAbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to create vars file '%s'", varsAbsPath)
		return
	}
	defer packerVarsFile.Close()

	if err = mergedTemplate.Execute(packerVarsFile, pv.Environment); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to execute merged template '%s'", mergedFileAbsPath)
		return
	}

	return
}

type Option func(*PackerVars)
