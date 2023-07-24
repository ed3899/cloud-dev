package binaries

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type Tool int

const (
	PackerID Tool = iota
	TerraformID
)

// type Terraform2I interface {
// 	Up() (err error)
// 	Down() (err error)
// }

// type Terraform2 struct {
// 	Path string
// 	Zip  *Zip
// }

// func (t *Terraform2) Up() (err error) {
// 	return
// }

// func (t *Terraform2) Down() (err error) {
// 	return
// }

type PackerAWSEnvironment struct {
	AWS_ACCESS_KEY                     string
	AWS_SECRET_KEY                     string
	AWS_IAM_PROFILE                    string
	AWS_USER_IDS                       []string
	AWS_AMI_NAME                       string
	AWS_INSTANCE_TYPE                  string
	AWS_REGION                         string
	AWS_EC2_AMI_NAME_FILTER            string
	AWS_EC2_AMI_ROOT_DEVICE_TYPE       string
	AWS_EC2_AMI_VIRTUALIZATION_TYPE    string
	AWS_EC2_AMI_OWNERS                 []string
	AWS_EC2_SSH_USERNAME               string
	AWS_EC2_INSTANCE_USERNAME          string
	AWS_EC2_INSTANCE_USERNAME_HOME     string
	AWS_EC2_INSTANCE_USERNAME_PASSWORD string
}

type PackerGeneralEnvironment struct {
	GIT_USERNAME                          string
	GIT_EMAIL                             string
	ANSIBLE_TAGS                          []string
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

type PackerEnvironment struct {
	AWS     *PackerAWSEnvironment
	General *PackerGeneralEnvironment
}

type Kind int

const (
	General Kind = iota
	AWS
)

type Template struct {
	Kind     Kind
	AbsPath  string
	Instance *template.Template
}

type VarsFile struct {
	Name        string
	Template    *Template
	AbsPath     string
	Environment any
}

func NewVarsFile(tool Tool, Kind Kind) (vf *VarsFile, err error) {
	const (
		generalVarsName     = "general.auto.tfvars"
		generalTemplateName = "GeneralPackerVars.tmpl"
		awsVarsName         = "aws_ami.auto.pkrvars.hcl"
		awsTemplateName     = "AWS_PackerVars.tmpl"
	)

	vf = &VarsFile{}
	switch tool {
		case PackerID:
			const (
				generalVarsName     = "general_ami.auto.pkrvars.hcl"
				generalTemplateName = "GeneralPackerVars.tmpl"
				awsVarsName         = "aws_ami.auto.pkrvars.hcl"
				awsTemplateName     = "AWS_PackerVars.tmpl"
			)
		case TerraformID:
			const (
				generalTerraformVarsName     = "general.auto.tfvars"
				generalTerraformTemplateName = "GeneralTerraformTfVars.tmpl"
				awsTerraformVarsName         = "aws.auto.tfvars"
				awsTerraformTemplateName     = "AWS_TerraformTfVars.tmpl"
			)
	}
}

func (vf *VarsFile) ParseTemplate(absPathToCloudDir string) {
	template.ParseFiles(vf.Template.AbsPath)
}

func (vf *VarsFile) CreateFile(absPathToCloudDir string) {

}

func (vf *VarsFile) ExecuteTemplate(absPathToCloudDir string) {

}

type Packer2I interface {
	Build() (err error)
}

type Packer2 struct {
	ID                  Tool
	Name                string
	AbsPathToExecutable string
	Cloud               string
	AbsPathToCloudDir   string
	AbsPathToPluginsDir string
	Zip                 *Zip
	VarsFiles           []*VarsFile
}

func NewPacker(cloud string) (packer *Packer2, err error) {
	const (
		name    = "packer"
		version = "1.9.1"
	)

	var (
		executableName = fmt.Sprintf("%s.exe", name)
		zipName        = fmt.Sprintf("%s.zip", name)
		os, arch       = utils.GetCurrentHostSpecs()
		url            = utils.CreateHashicorpURL(name, version, os, arch)
		depDirName     = utils.GetDependenciesDirName()
	)

	absPathToExecutable, err := filepath.Abs(filepath.Join(depDirName, name, executableName))
	if err != nil {
		err = errors.Wrapf(err, "failed to create binary path to: %s", executableName)
		return
	}
	absPathToCloudDir, err := filepath.Abs(filepath.Join(name, cloud))
	if err != nil {
		err = errors.Wrapf(err, "failed to create path to: %s", cloud)
		return
	}
	absPathToPluginsDir := filepath.Join(absPathToCloudDir, "plugins")
	zipPath, err := filepath.Abs(filepath.Join(depDirName, name, zipName))
	if err != nil {
		err = errors.Wrapf(err, "failed to craft zip path to: %s", zipName)
		return
	}
	contentLength, err := utils.GetContentLength(url)
	if err != nil {
		err = errors.Wrapf(err, "failed to get content length for: %s", url)
		return
	}

	packer = &Packer2{
		ID:                  PackerID,
		AbsPathToExecutable: absPathToExecutable,
		AbsPathToPluginsDir: absPathToPluginsDir,
		Cloud:               cloud,
		Zip: &Zip{
			Name:          zipName,
			Path:          zipPath,
			URL:           url,
			ContentLength: contentLength,
		},
	}

	return
}

func (p *Packer2) Init() (err error) {
	// Set PACKER_PLUGIN_PATH environment variable
	err = os.Setenv("PACKER_PLUGIN_PATH", p.AbsPathToPluginsDir)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while setting PACKER_PLUGIN_PATH environment variable")
		return
	}
	defer os.Unsetenv("PACKER_PLUGIN_PATH")

	// Get initial working directory
	initialLocation, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return
	}
	defer os.Chdir(initialLocation)

	//	Change working directory to where the packer files are depending on the cloud
	runLocation := filepath.Dir(p.AbsPathToExecutable)
	err = os.Chdir(runLocation)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while changing directory to run location")
		return
	}

	// Initialize
	cmd := exec.Command(p.AbsPathToExecutable, "init", ".")
	if cmdErr := utils.AttachCliToProcess(cmd); cmdErr != nil {
		err = errors.Wrapf(cmdErr, "Error occured while initializing %s for %s", p.Name, p.Cloud)
		return
	}
	return
}

func (p *Packer2) Build(cloud string) (err error) {

	return
}

func InitializeBinary() {

}
