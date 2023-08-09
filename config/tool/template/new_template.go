package template

import (
	"path/filepath"

	"github.com/ed3899/kumo/config/tool"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/utils/file"
	"github.com/samber/oops"
)

func NewTemplate(
	options ...Option,
) (
	template *Template,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("NewTemplate").
				With("options", options)

		option Option
	)

	template = &Template{}
	for _, option = range options {
		if err = option(template); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Failed to execute option")
			return
		}
	}

	return
}

func WithAbsPath(
	kumoExecAbsPath string,
) (
	option Option,
) {

	option = func(template *Template) (err error) {
		template.AbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.MERGED_TEMPLATE,
		)

		return
	}

	return
}

func (t *Template) Create(
	file_MergedFilesTo file.MergeFilesToF,
	toolConfig tool.ToolConfig,
) (err error) {
	var (
		oopsBuilder = oops.
			Code("Merge").
			With("template", t).
			With("tool", toolConfig)
	)

	if err = file_MergedFilesTo(
		t.AbsPath,
		toolConfig.AbsPathTo.TemplateFile.General,
		toolConfig.AbsPathTo.TemplateFile.Cloud,
	); err != nil {
		err = oopsBuilder.
			Wrapf(
				err,
				"Failed to merge general and cloud template '%s' and '%s'",
				toolConfig.AbsPathTo.TemplateFile.General,
				toolConfig.AbsPathTo.TemplateFile.Cloud,
			)

		return
	}

	return
}

func (t *Template) Remove(os_Remove func(string) error) (err error) {
	var (
		oopsBuilder = oops.
			Code("CallRemove")
	)

	if err = os_Remove(t.AbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to remove template '%s'", t.AbsPath)

		return
	}

	return
}

type Template struct {
	AbsPath string
}

type Option func(*Template) error
