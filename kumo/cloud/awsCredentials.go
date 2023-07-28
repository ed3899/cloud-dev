package cloud

type AwsCredentials struct {
	AccessKeyId     string
	SecretAccessKey string
}

func (ac *AwsCredentials) Set() (err error) {
	return
}

func (ac *AwsCredentials) Unset() (err error) {
	return
}
