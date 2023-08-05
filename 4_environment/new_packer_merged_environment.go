package environment

type PackerMergedEnvironment struct {
	General PackerGeneralEnvironmentI
	Cloud   PackerCloudEnvironmentI
}

func NewPackerMergedEnvironment(
	NewPackerGeneralEnvironment NewPackerGeneralEnvironmentF,
	NewPackerCloudEnvironment NewPackerCloudEnvironmentF,
) (
	mergedEnvironment PackerMergedEnvironment,
) {

	mergedEnvironment = PackerMergedEnvironment{
		General: NewPackerGeneralEnvironment(),
		Cloud:   NewPackerCloudEnvironment(),
	}

	return

}
