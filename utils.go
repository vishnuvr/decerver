package decerver

import (
	"fmt"
	"os"
	"path"
)

func AbsolutePath(Datadir string, filename string) string {
	if path.IsAbs(filename) {
		return filename
	}
	return path.Join(Datadir, filename)
}

func initDir(Datadir string) error {
	_, err := os.Stat(Datadir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Data directory '%s' doesn't exist, creating it\n", Datadir)
			mdaErr := os.MkdirAll(Datadir, 0777)
			if mdaErr != nil {
				return mdaErr
			}
		}
	} else {
		return err
	}
	return nil
}
