package packer

type PackerI interface {
	Init() (err error)
	Build() (err error)
}