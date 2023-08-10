package templates

func CloudTemplateName(
	args *CloudTemplateNameArgs,
) string {
	return args.fmt_Sprintf("%s.tmpl", args.CloudName())
}

type CloudTemplateNameF func(*CloudTemplateNameArgs) string

type CloudTemplateNameArgs struct {
	CloudName   func() string
	fmt_Sprintf func(string, ...any) string
}
