package template

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/constants"
	"github.com/ed3899/kumo/tool_manager"
	"github.com/ed3899/kumo/utils/file"
	"github.com/samber/oops"
)

func NewMergedTemplate(
	options ...Option,
) (
	template *MergedTemplate,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("NewTemplate").
				With("options", options)

		option Option
	)

	template = &MergedTemplate{}
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
	tool tool_manager.ToolManager,
	mergedFilesTo file.MergeFilesToF,
	kumoExecAbsPath string,
) (option Option) {
	var (
		oopsBuilder = oops.
				Code("WithInstance").
				With("tool", tool)

		mergedTemplateAbsPath = filepath.Join(
			kumoExecAbsPath,
			constants.TEMPLATES_DIR_NAME,
			constants.MERGED_TEMPLATE,
		)
	)

	option = func(template *MergedTemplate) (err error) {
		if err = mergedFilesTo(
			mergedTemplateAbsPath,
			tool.AbsPathTo.TemplateFile.General,
			tool.AbsPathTo.TemplateFile.Cloud,
		); err != nil {
			err = oopsBuilder.
				Wrapf(
					err,
					"Failed to merge general and cloud template '%s' and '%s'",
					tool.AbsPathTo.TemplateFile.General,
					tool.AbsPathTo.TemplateFile.Cloud,
				)

			return
		}

		if template.Instance, err = template.Instance.ParseFiles(mergedTemplateAbsPath); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Failed to parse merged template '%s'", mergedTemplateAbsPath)
			return
		}
	}

	return
}

type MergedTemplate struct {
	Instance *template.Template
}

type Option func(*MergedTemplate) error
