package binaries

type Terraform2I interface {
	Init(cloud Cloud) (err error)
	Up(cloud Cloud) (err error)
	Destroy(cloud Cloud) (err error)
}

type Terraform2 struct {
}
