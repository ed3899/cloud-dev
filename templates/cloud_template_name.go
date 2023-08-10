package templates

import (
	"github.com/ed3899/kumo/common/functions"
)

func CloudTemplateNameWith(
	args *functions.CloudTemplateNameArgs,
) functions.CloudTemplateName {
	return func() string {
		return args.Fmt_Sprintf("%s.tmpl", args.CloudName())
	}
}
