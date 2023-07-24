package utils

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

type Tool int

const (
	Packer Tool = iota
	Terraform
)

type Kind int

const (
	General Kind = iota
	AWS
)

func PickTool(tool Tool) (string, error) {
	switch tool {
	case Packer:
		return "packer", nil
	case Terraform:
		return "terraform", nil
	default:
		return "", errors.New("Invalid tool")
	}
}

func PickKind(kind Kind) (string, error) {
	switch kind {
	case General:
		return "general", nil
	case AWS:
		return "aws", nil
	default:
		return "", errors.New("Invalid kind")
	}
}

func PickName(tool Tool, kind Kind) (string, error) {
	switch tool {
	case Packer:
		switch kind {
		case General:
			return "GeneralPackerVars.tmpl", nil
		case AWS:
			return "AWS_PackerVars.tmpl", nil
		default:
			return "", errors.New("Invalid kind")
		}
	case Terraform:
		switch kind {
		case General:
			return "GeneralTerraformTfVars.tmpl", nil
		case AWS:
			return "AWS_TerraformTfVars.tmpl", nil
		default:
			return "", errors.New("Invalid kind")
		}
	default:
		return "", errors.New("Invalid tool")
	}
}

type ComputedProps struct {
	Tool string
	Kind string
	Name string
}

func PickProps(props *TemplateProps) (computedProps *ComputedProps, err error) {
	tool, err := PickTool(props.Tool)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while picking tool for template")
		return nil, err
	}

	kind, err := PickKind(props.Kind)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while picking kind for template")
		return nil, err
	}

	name, err := PickName(props.Tool, props.Kind)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while picking name for template")
		return nil, err
	}

	computedProps = &ComputedProps{
		Tool: tool,
		Kind: kind,
		Name: name,
	}

	return
}

type TemplateProps struct {
	Tool Tool
	Kind Kind
}

func CreateTemplate(props *TemplateProps) (tmpl *template.Template, err error) {
	computedProps, err := PickProps(props)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while picking props for template")
		return nil, err
	}

	templatePath, err := filepath.Abs(filepath.Join("templates", computedProps.Tool, computedProps.Kind, computedProps.Name))
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s template file", computedProps.Name)
		return nil, err
	}
	tmpl, err = template.ParseFiles(templatePath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while parsing template for %s", templatePath)
		return nil, err
	}

	return
}

func CreateVarsFile(props *TemplateProps) (varsFile *os.File, err error) {
	computedProps, err := PickProps(props)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while picking props for template")
		return nil, err
	}

	varsFileAbsPath, err := filepath.Abs(filepath.Join(computedProps.Tool, computedProps.Kind, computedProps.Name))
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s file", computedProps.Name)
		return nil, err
	}

	varsFile, err = os.Create(varsFileAbsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while creating %s file", computedProps.Name)
		return nil, err
	}
	varsFile.Close()

	return
}

type ExecuteTemplateProps struct {
	Template        *template.Template
	VarsFileAbsPath string
	Environment     any
}

func ExecuteTemplate(props *ExecuteTemplateProps) (err error) {
	tmpl := props.Template

	varsFile, err := os.Open(props.VarsFileAbsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while opening %s file", props.VarsFileAbsPath)
		return err
	}
	defer varsFile.Close()

	err = tmpl.Execute(varsFile, props.Environment)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while executing %s template file", tmpl.Name())
		return err
	}
	return nil
}

func ParseTemplate(absPathToTmpl string) (template *template.Template, err error) {
	template, err = template.ParseFiles(absPathToTmpl)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while parsing template for %s", absPathToTmpl)
		return
	}
	return
}

func CreateFile(absPathToFile string) (file *os.File, err error) {
	file, err = os.Create(absPathToFile)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while creating %s file", absPathToFile)
		return
	}
	defer file.Close()
	return
}

func ()  {
	
}