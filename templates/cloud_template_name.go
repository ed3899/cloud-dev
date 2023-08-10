package templates

func CloudTemplateName(
	CloudName func() string,
	fmt_Sprintf func(string, ...any) string,
) string {
	return fmt_Sprintf("%s.tmpl", CloudName())
}
