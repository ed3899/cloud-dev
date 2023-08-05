package environment

type TerraformGeneralRequired struct {
	ALLOWED_IP string
}

type TerraformGeneralEnvironment struct {
	Required TerraformGeneralRequired
}

func (tae TerraformGeneralEnvironment) IsEnvironment() bool {
	return true
}

func (tae TerraformGeneralEnvironment) IsGeneralEnvironment() bool {
	return true
}