package workflows

type PackerVars struct {
	Name 	 string
	AbsPath  string
}

type TerraformVars struct {

}

type HashicorpVars2 struct {
	Name     string
	AbsPath  string
	Template *Template2
}

func NewHashicorpVars2(template *Template2) (hv *HashicorpVars2, err error) {

	return
}
