package workflows


import (
	"github.com/ed3899/kumo/binaries/instances"
	"github.com/samber/oops"
)

func Build() (err error) {
	var (
		oopsBuilder = oops.
			Code("build_failed")
		packer *instances.Packer
	)

	// 1. Instantiate Packer
	if packer, err = instances.NewPacker(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while instantiating Packer")
		return
	}
	// 2. Download and install if needed
	if packer.IsNotInstalled() {
		
	}

	// 3. CloudSetup

	// 4. ToolSetup

	// 5. Create template

	// 6. Create hashicorp vars

	// 7. Change to right directory and defer change back

	// 8. Set plugin path

	// 9. Initialize

	// 10. Build

	return
}
