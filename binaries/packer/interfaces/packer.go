package packer

type Packer interface {
	Init() (err error)
	Build() (err error)
}