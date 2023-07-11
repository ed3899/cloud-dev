package utils

import (
	"os"

	"github.com/pkg/errors"
)

func DirExist(dirpath string) (bool, error) {
	_, err := os.Stat(dirpath)

	switch {
	case os.IsNotExist(err):
		return false, nil
	case os.IsExist(err):
		return true, nil
	case err != nil:
		err = errors.Wrap(err, "Error occurred while checking existence of directory")
		return false, err
	default:
		return true, nil
	}
}

func DirNotExist(dirpath string) (bool, error) {
	e, err := DirExist(dirpath)
	if err != nil {
		return false, nil
	}
	return !e, nil
}
