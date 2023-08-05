package environment

type PackerEnvironment struct {
	General PackerGeneralEnvironmentI
	Cloud   PackerCloudEnvironmentI
}

func NewPackerEnvironment(
	NewPackerGeneralEnvironment NewPackerGeneralEnvironmentF,
	NewPackerCloudEnvironment NewPackerCloudEnvironmentF,
) (
	packerEnvironment PackerEnvironment,
) {

	packerEnvironment = PackerEnvironment{
		General: NewPackerGeneralEnvironment(),
		Cloud:   NewPackerCloudEnvironment(),
	}

	return

}
