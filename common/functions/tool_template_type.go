package functions

type ToolTemplate func(*ToolTemplateArgs) string

type ToolTemplateArgs struct {
	CloudTemplateName string
	Fmt_Sprintf       func(string, ...any) string
}
