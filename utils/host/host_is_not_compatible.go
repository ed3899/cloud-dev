package host

func HostIsNotCompatible() (notCompatible bool) {
	return !HostIsCompatible()
}