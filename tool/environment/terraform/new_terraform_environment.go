package environment

type TerraformEnvironment struct {
	General TerraformGeneralEnvironmentI
	Cloud   TerraformCloudEnvironmentI
}

func NewTerraformEnvironment(
	NewTerraformGeneralEnvironment TerraformGeneralEnvironmentI,
	NewTerraformCloudEnvironment TerraformCloudEnvironmentI,
) (terraformEnvironment TerraformEnvironment) {

	terraformEnvironment = TerraformEnvironment{
		General: NewTerraformGeneralEnvironment,
		Cloud:   NewTerraformCloudEnvironment,
	}

	return
}
