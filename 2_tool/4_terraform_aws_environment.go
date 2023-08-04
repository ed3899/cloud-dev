package tool

type TerraformAwsRequired struct {
	AWS_REGION        string
	AWS_INSTANCE_TYPE string
	AMI_ID            string
	KEY_NAME          string
	SSH_PORT          int
	IP_FILE_NAME      string
	USERNAME          string
}

type TerraformAwsOptional struct {
	AWS_EC2_INSTANCE_VOLUME_TYPE string
	AWS_EC2_INSTANCE_VOLUME_SIZE int
}

type TerraformAwsEnvironment struct {
	Required TerraformAwsRequired
	Optional TerraformAwsOptional
}
