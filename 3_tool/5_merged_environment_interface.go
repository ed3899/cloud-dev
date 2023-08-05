package tool

type Merged[E PackerAwsEnvironment] struct {
	General E
	Cloud   E
}
