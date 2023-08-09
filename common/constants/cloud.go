package constants

type CloudKind int

const (
	Aws CloudKind = iota
)

func (c CloudKind) String() (s string) {
	switch c {
	case Aws:
		s = "aws"
	}

	return
}
