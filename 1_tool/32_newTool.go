package tool

type Data struct {
	absPath     string
	environment interfaces.Environment
}

type Templates struct {
	general *Data
	cloud   *Data
}

type Tool struct {
	kind              constants.Kind
	name              string
	version           string
	executableAbsPath string
	runDir            string
	pluginDir         string
	templates         *Templates
}

func NewTool()  {

}