package binz

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
)

type KumoConfigContent struct {
	AWS_ACCESS_KEY                        string
	AWS_SECRET_KEY                        string
	AWS_IAM_PROFILE                       string
	AWS_USER_IDS                          []string
	AWS_AMI_NAME                          string
	AWS_INSTANCE_TYPE                     string
	AWS_REGION                            string
	AWS_EC2_AMI_NAME_FILTER               string
	AWS_EC2_AMI_ROOT_DEVICE_TYPE          string
	AWS_EC2_AMI_VIRTUALIZATION_TYPE       string
	AWS_EC2_AMI_OWNERS                    []string
	AWS_EC2_SSH_USERNAME                  string
	AWS_EC2_INSTANCE_USERNAME             string
	AWS_EC2_INSTANCE_USERNAME_HOME        string
	AWS_EC2_INSTANCE_USERNAME_PASSWORD    string
	AWS_EC2_INSTANCE_SSH_KEY_NAME         string
	AWS_EC2_INSTANCE_VOLUME_TYPE          string
	AWS_EC2_INSTANCE_VOLUME_SIZE          string
	GIT_USERNAME                          string
	GIT_EMAIL                             string
	TOOLS                                 string
	GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC string
	PULUMI_PERSONAL_ACCESS_TOKEN          string
}

type KumoConfig struct {
	YAML_AbsPath string
	ParsedEnv    string
	Content      *KumoConfigContent
}

func GetKumoConfig() (kc *KumoConfig, err error) {
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
			kc.YAML_AbsPath = path
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
	if utils.FilePresent(kc.YAML_AbsPath) {
		return kc, nil
	}

	return kc, errors.New("kumo config file not found")
}

func TransformKumoConfig(kc *KumoConfig) (err error) {
	// Parse the yaml file
	
}
