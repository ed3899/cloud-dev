package tool

func ToolTemplateWith[
	Formatter ~func(string, ...any) string,
	CloudName ~func() string,
](
	fmt_Sprintf Formatter,
	cloudName CloudName,
) ToolTemplate {
	return func() string {
		return fmt_Sprintf("%s.tmpl", cloudName())
	}
}

type ToolTemplate func() string
