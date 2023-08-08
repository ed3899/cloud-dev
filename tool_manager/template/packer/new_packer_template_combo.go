package packer

import (
	"path/filepath"

	"github.com/ed3899/kumo/cloud"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/tool_manager"
	"github.com/ed3899/kumo/utils/file"
	"github.com/samber/oops"
)

func NewMergedTemplate(
	options ...Option,
) (
	packerMergedTemplate *MergedTemplate,
) {
	var (
		option Option
	)

	packerMergedTemplate = &MergedTemplate{}
	for _, option = range options {
		option(packerMergedTemplate)
	}

	return
}

func WithAbsPathFor(kumoExecAbsPath string, cloud cloud.Cloud, tool tool_manager.ToolManager) (option Option) {
	option = func(packerTemplate *MergedTemplate) (err error) {
		packerTemplate.AbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.PACKER_TEMP,
		)

		return
	}

	return
}

func WithTemplatesFor(
	cloud cloud.Cloud,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
			Code("WithCloudChoice").
			With("cloud", cloud).
			With("kumoExecAbsPath", kumoExecAbsPath)
	)

	option = func(packerTemplate *MergedTemplate) (err error) {
		packerTemplate.GeneralTemplateFileAbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.PACKER,
			constants.GENERAL_DIR_NAME,
			constants.PACKER_GENERAL_VARS_TEMPLATE,
		)

		switch cloud.Kind {
		case constants.Aws:
			packerTemplate.CloudTemplateFileAbsPath = filepath.Join(
				kumoExecAbsPath,
				constants.TEMPLATES_DIR_NAME,
				constants.PACKER,
				constants.AWS,
				constants.PACKER_AWS_VARS_TEMPLATE,
			)
		default:
			err = oopsBuilder.
				Wrapf(err, "Unsupported cloud '%s'", cloud.Name)
			return
		}

		return
	}

	return
}

func (ptc *MergedTemplate) Create(fileMerger file.MergeFilesToF) (err error) {
	var (
		oopsBuilder = oops.
			Code("OutputMerge").
			With("fileMerger", fileMerger)
	)

	if err = fileMerger(ptc.AbsPath, ptc.GeneralTemplateFileAbsPath, ptc.CloudTemplateFileAbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to merge general and cloud template '%s' and '%s'", ptc.GeneralTemplateFileAbsPath, ptc.CloudTemplateFileAbsPath)
		return
	}

	return
}

func (ptc *MergedTemplate) Remove(fileRemover func(name string) error) (err error) {
	var (
		oopsBuilder = oops.
			Code("RemoveMerge").
			With("fileRemover", fileRemover)
	)

	if err = fileRemover(ptc.AbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to remove merged template '%s'", ptc.AbsPath)
		return
	}

	return
}

type MergedTemplate struct {
	AbsPath                    string
	GeneralTemplateFileAbsPath string
	CloudTemplateFileAbsPath   string
}

type Option func(*MergedTemplate) error
