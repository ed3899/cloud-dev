package tool

func ToolTemplate(
	CloudName func() string,
	fmt_Sprintf func(string, ...any) string,
) string {
	return fmt_Sprintf("%s.tmpl", CloudName())
}

type ToolTemplateF func(*ToolTemplateArgs) string

type ToolTemplateArgs struct {
	CloudTemplateName func(func() string, func(string, ...any) string) string
}

func T(
	ToolTemplate *ToolTemplateArgs,
) {}
