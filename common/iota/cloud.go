package iota

import (
	"log"

	"github.com/samber/oops"
)

type Cloud int

const (
	Aws Cloud = iota
)

func (c Cloud) Name() string {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Cloud").
		Code("Name")

	switch c {
	case Aws:
		return "aws"

	default:
		err := oopsBuilder.
			Errorf("unknown cloud: %#v", c)

		log.Fatalf("%+v", err)

		return ""
	}
}

func (c Cloud) Template() Template {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Cloud").
		Code("Templates")

	switch c {
	case Aws:
		return Template{
			cloud: "aws.tmpl",
			base:  "base.tmpl",
		}

	default:
		err := oopsBuilder.
			Errorf("unknown cloud: %#v", c)

		log.Fatalf("%+v", err)

		return Template{}
	}
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
