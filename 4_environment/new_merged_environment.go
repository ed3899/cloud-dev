package environment

type MergedEnvironment struct {
	General GeneralEnvironmentI
	Cloud   CloudEnvironmentI
}

func NewMergedEnvironment(
	generalEnvironment GeneralEnvironmentI,
	packerCloudEnvironment CloudEnvironmentI,
) (
	mergedEnvironment MergedEnvironment,
) {

	mergedEnvironment = MergedEnvironment{
		General: generalEnvironment,
		Cloud:   packerCloudEnvironment,
	}

	return

}
