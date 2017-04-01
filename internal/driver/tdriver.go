package driver

import "sync"

// TDriver is driver impl for test.
// Don't use expect in test.
type TDriver struct {
	CountOfExec int
	Versions    []uint64
	Created     bool
	JustCreated bool
	locker      sync.Mutex
}

// Close ...
func (td *TDriver) Close() error {
	return nil
}

// CreateVersionTableIfNotExists ...
func (td *TDriver) CreateVersionTableIfNotExists() error {
	if !td.Created {
		td.Created = true
		td.JustCreated = true
	}
	return nil
}

// Exec ...
func (td *TDriver) Exec(query string) error {
	td.CountOfExec++
	return nil
}

// GetCurrentVersion ...
func (td *TDriver) GetCurrentVersion() (uint64, error) {
	if len(td.Versions) == 0 {
		return 0, nil
	}
	return td.Versions[len(td.Versions)-1], nil
}

// Open ...
func (td *TDriver) Open(dsn string) error {
	return nil
}

// SetVersion ...
func (td *TDriver) SetVersion(version uint64) error {
	td.Versions = append(td.Versions, version)
	return nil
}

// DeleteVersion ...
func (td *TDriver) DeleteVersion(version uint64) error {
	var vs []uint64
	for _, v := range td.Versions {
		if v != version {
			vs = append(vs, v)
		}
	}
	td.Versions = vs
	return nil
}

// Transaction ...
func (td *TDriver) Transaction(fn func() error) error {
	td.locker.Lock()
	defer td.locker.Unlock()

	if err := fn(); err != nil {
		return err
	}
	return nil
}
