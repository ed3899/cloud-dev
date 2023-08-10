package tool

import (
	"github.com/ed3899/kumo/common/functions"
)

func ToolTemplateWith(
	args *functions.ToolTemplateArgs,
) functions.ToolTemplate {
	return func() string {
		return args.Fmt_Sprintf("%s.tmpl", args.CloudTemplateName)
	}
}
