package utils

import (
	"os"

	"github.com/sicozz/papyrus/utils/constants"
)

func InitFsDir() (err error) {
	_, err = os.Stat(constants.PathFsDir)
	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		err = os.MkdirAll(constants.PathFsDir, 0755)
	}

	return
}
