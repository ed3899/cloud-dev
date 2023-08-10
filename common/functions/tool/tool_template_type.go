package tool

type ToolTemplateWith func(*ToolTemplateArgs) string

type ToolTemplate func() string

type ToolTemplateArgs struct {
	CloudTemplateName string
	Fmt_Sprintf       func(string, ...any) string
}
