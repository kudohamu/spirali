package spirali

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Readable represents abstraction that readable binary from path.
type Readable interface {
	Read(path string) ([]byte, error)
}

// Dir represents directory.
type Dir struct {
	basePath string
}

// NewReadableFromDir ...
func NewReadableFromDir(basePath string) *Dir {
	return &Dir{
		basePath: basePath,
	}
}

// Read ...
func (d *Dir) Read(path string) ([]byte, error) {
	file, err := os.Open(filepath.Join(d.basePath, path))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Bindata represents just `bindata`.
type Bindata struct {
	asset    func(string) ([]byte, error)
	basePath string
}

// NewReadableFromBindata ...
func NewReadableFromBindata(asset func(string) ([]byte, error)) *Bindata {
	return &Bindata{
		asset: asset,
	}
}

// Read ...
func (b *Bindata) Read(path string) ([]byte, error) {
	return b.asset(filepath.Join(b.basePath, path))
}

// WithBasePath sets a base path when read data from bindata.
func (b *Bindata) WithBasePath(basePath string) *Bindata {
	b.basePath = basePath
	return b
}
