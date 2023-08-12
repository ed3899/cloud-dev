package iota

import (
	"github.com/samber/oops"
)

type Cloud int

const (
	Aws Cloud = iota
)

func (c Cloud) Iota() Cloud {
	return c
}

func (c Cloud) Name() (name string) {
	oops.
		In("common").
		In("iota").
		Tags("Cloud").
		Code("Name").Recoverf(
		func() {
			switch c {
			case Aws:
				name = "aws"

			default:
				panic(c)
			}
		},
		"unknown cloud",
	)

	return
}

func (c Cloud) Template() (template Template) {
	oops.
		In("common").
		In("iota").
		Tags("Cloud").
		Code("Templates").
		Recoverf(
			func() {
				switch c {
				case Aws:
					template = Template{
						cloud: "aws.tmpl",
						base:  "base.tmpl",
					}

				default:
					panic(c)
				}
			},
			"unknown cloud",
		)

	return
}

type Template struct {
	cloud string
	base  string
}

func (t Template) Cloud() string {
	return t.cloud
}

func (t Template) Base() string {
	return t.base
}
