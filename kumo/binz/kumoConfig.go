package binz

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

type KumoConfig struct {
	Path string
	Content      *KumoConfigContent
	ParsedEnv    string
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
			kc.Path = path
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

func ParseKumoConfig(kc *KumoConfig, kind string) (err error) {
	// Open the file
	ykccf, err := os.Open(kc.Path)
	if err != nil {
		return errors.Wrap(err, "failed to open kumo config file")
	}
	defer ykccf.Close()

	kcc := KumoConfigContent{}
	err = yaml.NewDecoder(ykccf).Decode(&kcc)
	if err != nil {
		return errors.Wrap(err, "failed to decode kumo config file")
	}

	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
