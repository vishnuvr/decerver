package files

import ()

type FileIO interface {
	Root() string
	Log() string
	Dapps() string
	Blockchains() string
	Filesystems() string
	Modules() string
	System() string
	InitPaths() error
	// Useful when you want to load a file inside of a directory gotten by the
	// 'Paths' object. Reads and returns the bytes.
	ReadFile(directory, name string) ([]byte, error)
	// Useful when you want to save a file into a directory gotten by the 'Paths'
	// object.
	WriteFile(directory, name string, data []byte) error
	// Useful when you want to load json encoded files into objects.
	UnmarshalJsonFromFile(directory, name string, object interface{}) error
	// Useful when you want to store json encoding of objects in files.
	MarshalJsonToFile(directory, name string, object interface{}) error
	// Convenience method for creating module directories.
	CreateModuleDirectory(moduleName string) error
	// Convenience method for creating directories.
	CreateDirectory(dir string) error
}
