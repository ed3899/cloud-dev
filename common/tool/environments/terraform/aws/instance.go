package aws

import "github.com/ed3899/kumo/common/utils"

type Required struct {
	AWS_REGION        string
	AWS_INSTANCE_TYPE string
	AMI_ID            string
	KEY_NAME          string
	SSH_PORT          int
	IP_FILE_NAME      string
	USERNAME          string
}

type Optional struct {
	AWS_EC2_INSTANCE_VOLUME_TYPE string
	AWS_EC2_INSTANCE_VOLUME_SIZE int
}

type Environment struct {
	Required *Required
	Optional *Optional
}

func (e *Environment) IsNotValidEnvironment() (isNotValidEnvironment bool) {
	return !utils.IsStructCompletellyFilled(e.Required)
}
