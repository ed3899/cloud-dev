package aws

type Environment struct {
	AWS_REGION                   string
	AWS_INSTANCE_TYPE            string
	AWS_EC2_INSTANCE_VOLUME_TYPE string
	AWS_EC2_INSTANCE_VOLUME_SIZE int
	AMI_ID                       string
}

func (e *Environment) IsNotValidEnvironment() (isNotValidEnvironment bool) {
	switch {
	case e.AWS_REGION == "":
		return true
	case e.AWS_INSTANCE_TYPE == "":
		return true
	case e.AMI_ID == "":
		return true
	default:
		return false
	}
}
