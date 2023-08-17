package iota

import (
	"log"

	"github.com/samber/oops"
)

type Cloud int

const (
	Aws Cloud = iota
)

func (c Cloud) Iota() Cloud {
	return c
}

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

func (c Cloud) Template() *Template {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Cloud").
		Code("Template")

	switch c {
	case Aws:
		return &Template{
			Cloud: "aws.tmpl",
			Base:  "base.tmpl",
		}

	default:
		err := oopsBuilder.
			Errorf("unknown cloud: %#v", c)

		log.Fatalf("%+v", err)

		return nil
	}
}

type Template struct {
	Cloud string
	Base  string
}
