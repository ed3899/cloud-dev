package tool

func ToolTemplateWith(
	fmtSprintf func(string, ...any) string,
) ToolTemplate {
	toolTemplate := func(cloudName string) string {
		return fmtSprintf("%s.tmpl", cloudName)
	}

	return toolTemplate
}

type ToolTemplate func(cloudName string) string
