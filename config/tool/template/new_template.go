package template

import (
	"path/filepath"
	text_template "text/template"

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

func WithInstance(
	toolConfig tool.ToolConfig,
	mergedFilesTo file.MergeFilesToF,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithInstance").
				With("tool", toolConfig)

		mergedTemplateAbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.MERGED_TEMPLATE,
		)
	)

	option = func(template *Template) (err error) {
		if err = mergedFilesTo(
			mergedTemplateAbsPath,
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

		if template.Instance, err = text_template.ParseFiles(mergedTemplateAbsPath); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Failed to parse merged template '%s'", mergedTemplateAbsPath)
			return
		}

		return
	}

	return
}

type Template struct {
	Instance *text_template.Template
}

type Option func(*Template) error
