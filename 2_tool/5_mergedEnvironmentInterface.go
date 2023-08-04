package tool

type Merged[E Environment] struct {
	General E
	Cloud   E
}