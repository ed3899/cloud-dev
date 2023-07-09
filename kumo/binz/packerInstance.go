package binz

type PackerI interface {
	Build()
}

type Packer struct {
	ExecutablePath string
}

func (p *Packer) Build() {

}