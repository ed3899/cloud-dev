package environment

type PackerEnvironment struct {
	General PackerGeneralEnvironmentI
	Cloud   PackerCloudEnvironmentI
}

func NewPackerEnvironment(
	NewPackerGeneralEnvironment NewPackerGeneralEnvironmentF,
	NewPackerCloudEnvironment NewPackerCloudEnvironmentF,
) (
	mergedEnvironment PackerEnvironment,
) {

	mergedEnvironment = PackerEnvironment{
		General: NewPackerGeneralEnvironment(),
		Cloud:   NewPackerCloudEnvironment(),
	}

	return

}
