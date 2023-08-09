package template

import (
	"path/filepath"
	text_template "text/template"

	"github.com/ed3899/kumo/config/tool"
	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/utils/file"
	"github.com/samber/oops"
)

func NewTemplateFile(
	options ...Option,
) (
	template *TemplateFile,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("NewTemplate").
				With("options", options)

		option Option
	)

	template = &TemplateFile{}
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

	option = func(template *TemplateFile) (err error) {
		template.AbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.MERGED_TEMPLATE,
		)

		return
	}

	return
}

func (t *TemplateFile) Create(
	file_MergedFilesTo file.MergeFilesToF,
	toolConfig tool.ToolConfig,
) (
	err error,
) {
	var (
		oopsBuilder = oops.
			Code("Merge").
			With("template", t).
			With("tool", toolConfig)
	)

	if err = file_MergedFilesTo(
		t.AbsPath,
		toolConfig.AbsPath.TemplateFile.General,
		toolConfig.AbsPath.TemplateFile.Cloud,
	); err != nil {
		err = oopsBuilder.
			Wrapf(
				err,
				"Failed to merge general and cloud template '%s' and '%s'",
				toolConfig.AbsPath.TemplateFile.General,
				toolConfig.AbsPath.TemplateFile.Cloud,
			)

		return
	}

	return
}

func (t *TemplateFile) TextTemplate(
	template_ParseFiles func(...string) (*text_template.Template, error),
) (
	textTemplate *text_template.Template,
	err error,
) {
	var (
		oopsBuilder = oops.
			Code("TextTemplate").
			With("template_ParseFiles", template_ParseFiles)
	)

	if textTemplate, err = template_ParseFiles(t.AbsPath); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Failed to parse template '%s'", t.AbsPath)

		return
	}

	return
}

func (t *TemplateFile) Remove(
	os_Remove func(string) error,
) (
	err error,
) {
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

type TemplateFile struct {
	AbsPath string
}

type Option func(*TemplateFile) error
