package workflows

import (
	"path/filepath"
	"text/template"

	"github.com/samber/oops"
	"github.com/spf13/viper"
)

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
	PackerGeneralEnvironment           *PackerGeneralEnvironment
}



type PackerAwsTemplate struct {
	Instance *template.Template
}

func NewPackerAwsTemplate() (packerAwsTemplate *PackerAwsTemplate, err error) {
	const (
		PACKER_AWS_TEMPLATE_NAME = "AWS_PackerVars.tmpl"
	)

	return
}

type AwsTemplateFile struct {
	Name    string
	AbsPath string
}

// type PackerMergeCombo struct {
// 	General *GeneralTemplateFile
// 	Aws     *AwsTemplateFile
// 	Merged  *MergedTemplateFile
// }

// func newPackerMergeCombo() (packerMergeCombo *PackerMergeCombo, err error) {
// 	var (
// 		oopsBuilder = oops.
// 			Code("new_packer_merge_combo_failed")
// 	)

// 	return
// }

// func NewPackerAwsTemplate() (packerAwsTemplate *PackerAwsTemplate, err error) {
// 	var (
// 		oopsBuilder = oops.
// 			Code("new_packer_aws_template_failed")
// 	)

// 	return

// }
