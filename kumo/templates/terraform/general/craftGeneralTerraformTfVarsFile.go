package templates

import (
	"os"
	"path/filepath"
	"text/template"

	templates_terraform_aws "github.com/ed3899/kumo/templates/terraform/aws"
	"github.com/ed3899/kumo/utils"
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

type VarsFile struct {
	Tool          Tool
	Kind          Kind
	ComputedProps *ComputedProps
	Template      *template.Template
	Instance      *os.File

	*templates_terraform_aws.AWS_TerraformEnvironment
	AMI_ID     string
	ALLOWED_IP string
}

func NewVars(tool Tool, kind Kind) (*VarsFile, error) {
	t := &VarsFile{
		Tool: tool,
		Kind: kind,
	}

	if err := t.setComputedTool(); err != nil {
		err = errors.Wrapf(err, "Error occurred while picking tool")
		return nil, err
	}

	if err := t.setComputedKind(); err != nil {
		err = errors.Wrapf(err, "Error occurred while picking kind")
		return nil, err
	}

	if err := t.setComputedName(); err != nil {
		err = errors.Wrapf(err, "Error occurred while picking name")
		return nil, err
	}

	if err := t.parseTemplate(); err != nil {
		err = errors.Wrapf(err, "Error occurred while parsing template")
		return nil, err
	}

	return t, nil
}

func (t *VarsFile) setComputedTool() (err error) {
	switch t.Tool {
	case Packer:
		t.ComputedProps.Tool = "packer"
	case Terraform:
		t.ComputedProps.Tool = "terraform"
	default:
		err = errors.New("Invalid tool")
	}
	return
}

func (t *VarsFile) setComputedKind() (err error) {
	switch t.Kind {
	case General:
		t.ComputedProps.Kind = "general"
	case AWS:
		t.ComputedProps.Kind = "aws"
	default:
		err = errors.New("Invalid kind")
	}
	return
}

func (t *VarsFile) setComputedName() (err error) {
	switch t.Tool {
	case Packer:
		switch t.Kind {
		case General:
			t.ComputedProps.Name = "GeneralPackerVars.tmpl"
		case AWS:
			t.ComputedProps.Name = "AWS_PackerVars.tmpl"
		default:
			err = errors.New("Invalid kind")
		}
	case Terraform:
		switch t.Kind {
		case General:
			t.ComputedProps.Name = "GeneralTerraformTfVars.tmpl"
		case AWS:
			t.ComputedProps.Name = "AWS_TerraformTfVars.tmpl"
		default:
			err = errors.New("Invalid kind")
		}
	default:
		err = errors.New("Invalid tool")
	}
	return
}

func (t *VarsFile) parseTemplate() (err error) {
	if t.ComputedProps.Tool == "" || t.ComputedProps.Kind == "" || t.ComputedProps.Name == "" {
		err = errors.New("Computed props not set")
		return err
	}

	templatePath, err := filepath.Abs(filepath.Join("templates", t.ComputedProps.Tool, t.ComputedProps.Kind, t.ComputedProps.Name))
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while crafting absolute path to %s template file", t.ComputedProps.Name)
		return err
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while parsing template for %s", templatePath)
		return err
	}

	t.Template = tmpl

	return nil
}

func (t *VarsFile) create() (varsFilePath string, err error) {

}

func CraftGeneralTerraformTfVarsFile(gte *VarsFile) (generalTerraformVarsPath string, err error) {

	// Create vars file
	varsFile, err := utils.CreateVarsFile(&utils.TemplateProps{
		Tool: utils.Terraform,
		Kind: utils.General,
	})
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while creating vars file")
		return "", err
	}

	err = utils.ExecuteTemplate(&utils.ExecuteTemplateProps{
		Template:        tmpl,
		VarsFileAbsPath: varsFile.Name(),
		Environment:     gte,
	})
	if err != nil {
		err = errors.Wrapf(err, "Error occurred while executing template")
		return "", err
	}

	generalTerraformVarsPath = varsFile.Name()

	return
}
