package deCerver

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

func InitDirs(Datadir string) {
	_, err := os.Stat(Datadir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Data directory '%s' doesn't exist, creating it\n", Datadir)
			os.MkdirAll(Datadir, 0777)
		}
	}
}