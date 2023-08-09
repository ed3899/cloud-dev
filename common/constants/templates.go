package constants

const (
	GENERAL_TEMPLATE = "general.tmpl"
	AWS_TEMPLATE     = "aws.tmpl"
	MERGED_TEMPLATE  = "merged.tmpl"
)

type TemplateKind int

const (
	General = iota
	Cloud
)

func (tk TemplateKind) General() (g string) {
	g = "general.tmpl"

	return
}

func (tk TemplateKind) Cloud() (c string) {
	
}
