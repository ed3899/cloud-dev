package iota

import "github.com/samber/oops"

type Cloud int

type awsC struct {

}


const (
	Base Cloud = iota
	Aws
)

func (c Cloud) Name() string {
	var choice string

	oops.
		In("iota").
		Tags("Cloud").
		Code("Name").
		Recoverf(func() {
			switch c {
			case Base:
				choice = "base"
			case Aws:
				choice = "aws"
			default:
				panic(c)
			}
		}, "Unknown cloud")

	return choice
}

func (c Cloud) Templates() (cloud string, base string) {
	oops.
		In("iota").
		Tags("Cloud").
		Code("TemplateName").
		Recoverf(func() {
			switch c {
			case Aws:
				cloud = "aws.tmpl"
			default:
				panic(c)
			}
		}, "Unknown cloud")

	base = "base.tmpl"

	return cloud, base
}
