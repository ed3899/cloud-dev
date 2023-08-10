package functions

type CloudTemplateNameWith func(*CloudTemplateNameArgs) CloudTemplateName

type CloudTemplateName func() string

type CloudTemplateNameArgs struct {
	CloudName   func() string
	Fmt_Sprintf func(string, ...any) string
}
