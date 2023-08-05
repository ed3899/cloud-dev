package tool

type TerraformGeneralRequired struct {
	ALLOWED_IP string
}

type TerraformGeneralEnvironment struct {
	Required TerraformGeneralRequired
}
