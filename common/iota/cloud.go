package iota

type Cloud int

const (
	General Cloud = iota
	Aws
)

func (c Cloud) Name() string {
	switch c {
	case General:
		return "general"
	case Aws:
		return "aws"
	default:
		panic("Unknown cloud")
	}
}

func (c Cloud) TemplateName() string {
	switch c {
	case General:
		return "general.tmpl"
	case Aws:
		return "aws.tmpl"
	default:
		panic("Unknown cloud")
	}
}
