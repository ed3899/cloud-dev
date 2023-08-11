package tool

func ToolTemplateNameWith(
	fmtSprintf func(string, ...any) string,
) ToolTemplateName {
	toolTemplateName := func(cloudName string) string {
		return fmtSprintf("%s.tmpl", cloudName)
	}

	return toolTemplateName
}

type ToolTemplateName func(cloudName string) string
