package binz

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/binz/download"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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

func (p *Packer) WritePackerAWS_VarsFile() (awsPackerHclPath string, err error) {
	envData := &AWS_PackerEnvironment{
		AWS_ACCESS_KEY:                     viper.GetString("AWS.AccessKeyId"),
		AWS_SECRET_KEY:                     viper.GetString("AWS.SecretAccessKey"),
		AWS_IAM_PROFILE:                    viper.GetString("AWS.IamProfile"),
		AWS_USER_IDS:                       viper.GetStringSlice("AWS.UserIds"),
		AWS_AMI_NAME:                       viper.GetString("AMI.Name"),
		AWS_INSTANCE_TYPE:                  viper.GetString("AWS.EC2.Instance.Type"),
		AWS_REGION:                         viper.GetString("AWS.Region"),
		AWS_EC2_AMI_NAME_FILTER:            viper.GetString("AMI.Base.Filter"),
		AWS_EC2_AMI_ROOT_DEVICE_TYPE:       viper.GetString("AMI.Base.RootDeviceType"),
		AWS_EC2_AMI_VIRTUALIZATION_TYPE:    viper.GetString("AMI.Base.VirtualizationType"),
		AWS_EC2_AMI_OWNERS:                 viper.GetStringSlice("AMI.Base.Owners"),
		AWS_EC2_SSH_USERNAME:               viper.GetString("AMI.Base.User"),
		AWS_EC2_INSTANCE_USERNAME:          viper.GetString("AMI.User"),
		AWS_EC2_INSTANCE_USERNAME_HOME:     viper.GetString("AMI.Home"),
		AWS_EC2_INSTANCE_USERNAME_PASSWORD: viper.GetString("AMI.Password"),
	}

	cwd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting current working directory")
		return "", err
	}

	awsTemplatePath := filepath.Join(cwd, "templates", "AWS_PackerVarsTemplate.txt")

	tmpl, err := template.ParseFiles(awsTemplatePath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while parsing Packer AWS Vars template file")
		return "", err
	}

	phcldir, err := utils.GetPackerHclDirPath()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting Packer HCL directory path")
		return "", err
	}

	awsPackerHclPath = filepath.Join(phcldir, "aws_ami.pkr.hcl")

	file, err := os.Create(awsPackerHclPath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while creating Packer AWS Vars file")
		return "", err
	}
	defer file.Close()

	err = tmpl.Execute(file, envData)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while executing Packer AWS Vars template file")
		return "", err
	}

	return awsPackerHclPath, nil

}

func (p *Packer) writePackerVarsFile(cloud string) (path string, err error) {
	switch cloud {
	case "aws":
		path, err = p.WritePackerAWS_VarsFile()
		if err != nil {
			err = errors.Wrap(err, "Error occurred while writing Packer AWS vars file")
			return "", err
		}
		return path, nil
	default:
		err = errors.Errorf("Cloud '%s' not supported", cloud)
		return "", err
	}
}

func (p *Packer) buildAMI_OnAWS() (err error) {
	phclfp, err := p.init()
	if err != nil {
		err = errors.Wrap(err, "Error occurred while initializing Packer with AWS config")
		return err
	}

	_, err = p.writePackerVarsFile("aws")
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
