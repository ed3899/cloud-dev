package interfaces

type PackerManifest interface {
	LastBuiltAmiId() (lastBuiltAmiId string)
}
