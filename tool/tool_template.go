package tool

import (
	"github.com/ed3899/kumo/common/functions"
)

func ToolTemplate(
	args *functions.ToolTemplateArgs,
) string {
	return args.Fmt_Sprintf("%s.tmpl", args.CloudTemplateName)
}


