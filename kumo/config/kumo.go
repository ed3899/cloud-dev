package config

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type AWS_Config struct {
	AccessKeyId     string
	SecretAccessKey string
	IAmProfile      string
	UserIds         []string
	Region          string
	EC2             struct {
		Instance struct {
			Type string
		}
		Volume struct {
			Type string
			Size int
		}
	}
}

type AMI_Config struct {
	Name string
	Base struct {
		Filter             string
		User               string
		RootDeviceType     string
		VirtualizationType string
		Owners             []string
	}
	User     string
	Home     string
	Password string
	Tools    []string
}

type GitConfig struct {
	Username string
	Email    string
}

type GitHubConfig struct {
	PersonalAccessTokenClassic string
}

type PulumiConfig struct {
	PersonalAccessToken string
}

type UpConfig struct {
	AMI_Id string
}

type KumoConfigContent struct {
	AWS    *AWS_Config
	AMI    *AMI_Config
	Git    *GitConfig
	GitHub *GitHubConfig
	Pulumi *PulumiConfig
	Up     *UpConfig
}

type KumoConfigI interface {
	GetKumoEnvironment() (*KumoEnvironment, error)
}

type KumoConfig struct {
	Path string
	Kind string
}

type KumoAWSEnvironment struct {
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
	AWS_EC2_INSTANCE_VOLUME_TYPE       string
	AWS_EC2_INSTANCE_VOLUME_SIZE       int
}

type KumoGeneralEnvironment struct {
	GIT_USERNAME string
	GIT_EMAIL    string
	ANSIBLE_TAGS []string
}

type KumoOptionalEnvironment struct {
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
	PULUMI_PERSONAL_ACCESS_TOKEN          string
}

type KumoEnvironment struct {
	*KumoAWSEnvironment
	*KumoGeneralEnvironment
	*KumoOptionalEnvironment
}

func GetKumoConfig(kind string) (kc *KumoConfig, err error) {
	// Get current working directory
	cwd, err := utils.GetCWD()
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		return nil, err
	}

	// Regex must match kumo.config.yaml or kumo.config.yml
	pattern := regexp.MustCompile(`^kumo\.config\.(yaml|yml)$`)

	// Walk the current working directory looking for the kumo config file
	// If found, set the path to the kumo config file
	// If not found, the path will remain empty string
	// This function can only prove existence of the file, not abscence!
	err = filepath.WalkDir(cwd, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			err = errors.Wrap(err, "failed to walk directory")
			return err
		}

		switch {
		case d.IsDir() && pattern.MatchString(d.Name()):
			return errors.New("found a directory but should be a file")
		case pattern.MatchString(d.Name()):
			log.Printf("Found kumo config file: %s", d.Name())
			kc.Path = path
			kc.Kind = kind
			return nil
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the kumo path exists. This is somehow obvious when getting
	// no error from the above function. However, it allows us to prove
	// the abscence of the file in case of an empty string passed as the path
	if utils.FilePresent(kc.Path) {
		return kc, nil
	}

	return kc, errors.New("kumo config file not found")
}

func (kc *KumoConfig) attachGeneralEnvironment(ke *KumoEnvironment, kcc *KumoConfigContent) {
	ke.KumoGeneralEnvironment.GIT_EMAIL = kcc.Git.Email
	ke.KumoGeneralEnvironment.GIT_USERNAME = kcc.Git.Username
	ke.KumoGeneralEnvironment.ANSIBLE_TAGS = kcc.AMI.Tools
}

func (kc *KumoConfig) attachOptionalEnvironment(ke *KumoEnvironment, kcc *KumoConfigContent) {
	ke.KumoOptionalEnvironment.GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC = kcc.GitHub.PersonalAccessTokenClassic
	ke.KumoOptionalEnvironment.PULUMI_PERSONAL_ACCESS_TOKEN = kcc.Pulumi.PersonalAccessToken
}

func (kc *KumoConfig) attachAWSKumoEnvironment(ke *KumoEnvironment, kcc *KumoConfigContent) {
	ke.KumoAWSEnvironment.AWS_ACCESS_KEY = kcc.AWS.AccessKeyId
	ke.KumoAWSEnvironment.AWS_SECRET_KEY = kcc.AWS.SecretAccessKey
	ke.KumoAWSEnvironment.AWS_IAM_PROFILE = kcc.AWS.IAmProfile
	ke.KumoAWSEnvironment.AWS_USER_IDS = kcc.AWS.UserIds
	ke.KumoAWSEnvironment.AWS_AMI_NAME = kcc.AMI.Name
	ke.KumoAWSEnvironment.AWS_INSTANCE_TYPE = kcc.AWS.EC2.Instance.Type
	ke.KumoAWSEnvironment.AWS_REGION = kcc.AWS.Region
	ke.KumoAWSEnvironment.AWS_EC2_AMI_NAME_FILTER = kcc.AMI.Base.Filter
	ke.KumoAWSEnvironment.AWS_EC2_AMI_ROOT_DEVICE_TYPE = kcc.AMI.Base.RootDeviceType
	ke.KumoAWSEnvironment.AWS_EC2_AMI_VIRTUALIZATION_TYPE = kcc.AMI.Base.VirtualizationType
	ke.KumoAWSEnvironment.AWS_EC2_AMI_OWNERS = kcc.AMI.Base.Owners
	ke.KumoAWSEnvironment.AWS_EC2_SSH_USERNAME = kcc.AMI.Base.User
	ke.KumoAWSEnvironment.AWS_EC2_INSTANCE_USERNAME = kcc.AMI.User
	ke.KumoAWSEnvironment.AWS_EC2_INSTANCE_USERNAME_HOME = kcc.AMI.Home
	ke.KumoAWSEnvironment.AWS_EC2_INSTANCE_USERNAME_PASSWORD = kcc.AMI.Password
	ke.KumoAWSEnvironment.AWS_EC2_INSTANCE_VOLUME_TYPE = kcc.AWS.EC2.Volume.Type
	ke.KumoAWSEnvironment.AWS_EC2_INSTANCE_VOLUME_SIZE = kcc.AWS.EC2.Volume.Size
}

func (kc *KumoConfig) GetKumoEnvironment() (ke *KumoEnvironment, err error) {
	// Open the file
	ykc, err := os.Open(kc.Path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open kumo config file")
	}
	defer ykc.Close()

	kcc := &KumoConfigContent{}
	err = yaml.NewDecoder(ykc).Decode(&kcc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode kumo config file")
	}

	ke = &KumoEnvironment{}
	kc.attachGeneralEnvironment(ke, kcc)
	kc.attachOptionalEnvironment(ke, kcc)

	switch kc.Kind {
	case "aws":
		kc.attachAWSKumoEnvironment(ke, kcc)
		return ke, nil
	default:
		return nil, errors.New("unknown kumo config kind")
	}
}
