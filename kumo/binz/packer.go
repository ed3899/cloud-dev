package binz

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/ed3899/kumo/binz/download"
	"github.com/ed3899/kumo/templates"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type AWS_PackerEnvironment struct {
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

type GeneralPackerEnvironment struct {
	*AWS_PackerEnvironment
	GIT_USERNAME                          string
	GIT_EMAIL                             string
	ANSIBLE_TAGS                          []string
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
}

type PackerI interface {
	GetPackerInstance(*download.Binaries) (*Packer, error)
	Build(string)
}

type Packer struct {
	ExecutablePath string
}

func (p *Packer) init() (phclfp string, err error) {
	phclfp, err = utils.GetPackerHclFilePath()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Packer HCL file path")
		return "", err
	}

	cmd := exec.Command(p.ExecutablePath, "init", phclfp)
	output, err := cmd.CombinedOutput()
	log.Print(string(output))
	if err != nil {
		err = errors.Wrap(err, "Error occurred while initializing Packer")
		return "", err
	}

	return phclfp, nil
}

func (p *Packer) buildAMI_OnAWS() (err error) {
	phclfp, err := p.init()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while initializing Packer with AWS config")
		return err
	}

	_, err = templates.CraftAWSPackerVarsFile()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while writing Packer AWS vars file")
		return err
	}

	cmd := exec.Command(p.ExecutablePath, "validate", phclfp)
	output, err := cmd.CombinedOutput()
	log.Print(string(output))
	if err != nil {
		err = errors.Wrap(err, "Error occurred while building AMI with AWS Config")
		return err
	}

	return nil
}

func (p *Packer) Build(cloud string) {
	switch cloud {
	case "aws":
		err := p.buildAMI_OnAWS()
		if err != nil {
			err = errors.Wrap(err, "Error occurred while building AMI on AWS")
			log.Fatal(err)
		}
	default:
		err := errors.Errorf("Cloud '%s' not supported", cloud)
		log.Fatal(err)
	}
}

func GetPackerInstance(bins *download.Binaries) (packer *Packer, err error) {
	// Create the absolute path to the executable
	ep := filepath.Join(bins.Packer.Dependency.ExtractionPath, "packer.exe")

	// Validate existence
	if utils.FileNotPresent(ep) {
		err = errors.New("Packer executable not found")
		return nil, err
	}

	packer = &Packer{
		ExecutablePath: ep,
	}

	return packer, nil
}
