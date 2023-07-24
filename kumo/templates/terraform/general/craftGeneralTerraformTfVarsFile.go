package templates

import (
	"os"
	"path/filepath"
	"text/template"

	templates_terraform_aws "github.com/ed3899/kumo/templates/terraform/aws"
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

type ComputedProps struct {
	Tool string
	Kind string
	Name string
}

type Environment struct {
	*templates_terraform_aws.AWS_TerraformEnvironment
	AMI_ID     string
	ALLOWED_IP string
}

type VarsFile struct {
	Tool          Tool
	Kind          Kind
	Environment   *Environment
	computedProps *ComputedProps
	template      *template.Template
	instance      *os.File
}

func NewVarsFile(tool Tool, kind Kind, env *Environment) (*VarsFile, error) {
	vf := &VarsFile{
		Tool: tool,
		Kind: kind,
		Environment: env,
	}

	if err := vf.setComputedTool(); err != nil {
		err = errors.Wrapf(err, "Error occurred while picking tool")
		return nil, err
	}

	if err := vf.setComputedKind(); err != nil {
		err = errors.Wrapf(err, "Error occurred while picking kind")
		return nil, err
	}

	if err := vf.setComputedName(); err != nil {
		err = errors.Wrapf(err, "Error occurred while picking name")
		return nil, err
	}

	if err := vf.parseTemplate(); err != nil {
		err = errors.Wrapf(err, "Error occurred while parsing template")
		return nil, err
	}

	if err := vf.create(); err != nil {
		err = errors.Wrapf(err, "Error occurred while creating vars file")
		return nil, err
	}

	if err := vf.executeTemplate(); err != nil {
		err = errors.Wrapf(err, "Error occurred while executing template")
		return nil, err
	}

	return vf, nil
}

func (vf *VarsFile) setComputedTool() (err error) {
	switch vf.Tool {
	case Packer:
		vf.computedProps.Tool = "packer"
	case Terraform:
		vf.computedProps.Tool = "terraform"
	default:
		err = errors.New("Invalid tool")
	}
	return
}

func (vf *VarsFile) setComputedKind() (err error) {
	switch vf.Kind {
	case General:
		vf.computedProps.Kind = "general"
	case AWS:
		vf.computedProps.Kind = "aws"
	default:
		err = errors.New("Invalid kind")
	}
	return
}

func (vf *VarsFile) setComputedName() (err error) {
	switch vf.Tool {
	case Packer:
		switch vf.Kind {
		case General:
			vf.computedProps.Name = "GeneralPackerVars.tmpl"
		case AWS:
			vf.computedProps.Name = "AWS_PackerVars.tmpl"
		default:
			err = errors.New("Invalid kind")
		}
	case Terraform:
		switch vf.Kind {
		case General:
			vf.computedProps.Name = "GeneralTerraformTfVars.tmpl"
		case AWS:
			vf.computedProps.Name = "AWS_TerraformTfVars.tmpl"
		default:
			err = errors.New("Invalid kind")
		}
	default:
		err = errors.New("Invalid tool")
	}
	return
}

func (vf *VarsFile) parseTemplate() (err error) {
	if vf.computedProps.Tool == "" || vf.computedProps.Kind == "" || vf.computedProps.Name == "" {
		err = errors.New("Computed props not set")
		return err
	}

	templatePath, err := filepath.Abs(filepath.Join("templates", vf.computedProps.Tool, vf.computedProps.Kind, vf.computedProps.Name))
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s template file", vf.computedProps.Name)
		return err
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while parsing template for %s", templatePath)
		return err
	}

	vf.template = tmpl

	return
}

func (vf *VarsFile) create() (err error) {
	if vf.computedProps.Tool == "" || vf.computedProps.Kind == "" || vf.computedProps.Name == "" {
		err = errors.New("Computed props not set")
		return err
	}

	varsFileAbsPath, err := filepath.Abs(filepath.Join(vf.computedProps.Tool, vf.computedProps.Kind, vf.computedProps.Name))
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s file", vf.computedProps.Name)
		return err
	}

	varsFile, err := os.Create(varsFileAbsPath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while creating %s file", vf.computedProps.Name)
		return err
	}
	varsFile.Close()

	vf.instance = varsFile

	return
}

func (t *VarsFile) executeTemplate() (err error) {
	if t.template == nil {
		err = errors.New("Template not set")
		return
	}

	varsFile, err := os.Open(t.instance.Name())
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while opening %s file", t.instance.Name())
		return
	}
	defer varsFile.Close()

	err = t.template.Execute(varsFile, t.Environment)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while executing %s template file", t.template.Name())
		return
	}
	return
}

func (t *VarsFile) GetPath() (absPath string, err error) {
	if t.instance == nil {
		err = errors.New("Instance not set")
		return
	}
	absPath = t.instance.Name()
	return
}
