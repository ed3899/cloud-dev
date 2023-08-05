package templates_parser

import (
	"text/template"

	"github.com/samber/oops"
)

func NewPackerAwsTemplate(kumoExecAbsPath string) (packerAwsTemplate template.Template, err error) {
	var (
		oopsBuilder = oops.
			Code("NewPackerAwsTemplate").
			With("kumoExecAbsPath", kumoExecAbsPath)
	)

	template.ParseFiles(kumoExecAbsPath,)

	return
}
