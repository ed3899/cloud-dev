package packer

func PackerUrl(
	args *PackerUrlArgs,
) (string ){

	return args.filepath_Join()


}

type PackerUrlArgs struct {
	filepath_Join func(...string) string
}