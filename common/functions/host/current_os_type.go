package host

type CurrentOs func() Os

type Os string

func (o Os) String() string {
	return string(o)
}
