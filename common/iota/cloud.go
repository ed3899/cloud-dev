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

	defer func() {
		if r := recover(); r != nil {
			err := oopsBuilder.Errorf("%v", r)
			log.Fatalf("panic: %+v", err)
		}
	}()

	switch c {
	case Aws:
		return "aws"

	default:
		panic(c)
	}
}

func (c Cloud) Template() Template {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Cloud").
		Code("Templates")

	defer func() {
		if r := recover(); r != nil {
			err := oopsBuilder.Errorf("%v", r)
			log.Fatalf("panic: %+v", err)
		}
	}()

	switch c {
	case Aws:
		return Template{
			cloud: "aws.tmpl",
			base:  "base.tmpl",
		}

	default:
		panic(c)
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
