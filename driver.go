package spirali

import "github.com/kudohamu/spirali/internal/driver"

// Driver is interface of database driver.
type Driver interface {
	Close() error
	CreateVersionTableIfNotExists() error
	Exec(query string) error
	GetCurrentVersion() (uint64, error)
	Open(dsn string) error
	SetVersion(version uint64) error
	Transaction(fn func() error) error
}

// NewDriver separates out actual sql driver.
func NewDriver(c *Config) (Driver, error) {
	switch c.Driver() {
	case "mysql":
		return &driver.Mysql{}, nil
	}
	return nil, ErrUnknownDriver
}
