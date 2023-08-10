package host

type CurrentArch func() Arch

type Arch string

func (a Arch) String() string {
	return string(a)
}
