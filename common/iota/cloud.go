package iota

import "github.com/samber/oops"

type Cloud int

const (
	Aws Cloud = iota
)

func (c Cloud) Name() (string, error) {
	oopsBuilder := oops.
		In("common").
		In("iota").
		Tags("Cloud").
		Code("Name")

	switch c {
	case Aws:
		return "aws", nil

	default:
		err := oopsBuilder.
			Errorf("unknown cloud: %#v", c)

		return "", err
	}
}

func (c Cloud) Template() (Template, error) {
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
		}, nil

	default:
		err := oopsBuilder.
			Errorf("unknown cloud: %#v", c)

		return Template{}, err
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
