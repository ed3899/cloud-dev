package iota

import "github.com/samber/oops"

type Cloud int

const (
	Base Cloud = iota
	Aws
)

func (c Cloud) Name() (name string) {
	oops.
		In("iota").
		Tags("Cloud").
		Code("Name").
		Recoverf(func() {
			switch c {
			case Base:
				name = "base"
			case Aws:
				name = "aws"
			default:
				panic(c)
			}
		}, "Unknown cloud")

	return
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
