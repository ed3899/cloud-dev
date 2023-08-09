package interfaces

type UrlGetter[U ~string] interface {
	Url() U
}