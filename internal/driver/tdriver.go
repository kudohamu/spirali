package driver

import (
	"sync"
	"time"
)

// TDriver is driver impl for test.
// Don't use expect in test.
type TDriver struct {
	CountOfExec int
	Rows        []*Row
	Created     bool
	JustCreated bool
	locker      sync.Mutex
}

// Row is TDriver's row
type Row struct {
	Version   uint64
	CreatedAt time.Time
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

// GetAppliedTimeList returns list of migration applied times.
func (td *TDriver) GetAppliedTimeList() (map[uint64]time.Time, error) {
	data := map[uint64]time.Time{}

	for _, r := range td.Rows {
		data[r.Version] = r.CreatedAt
	}

	return data, nil
}

// GetCurrentVersion ...
func (td *TDriver) GetCurrentVersion() (uint64, error) {
	if len(td.Rows) == 0 {
		return 0, nil
	}
	return td.Rows[len(td.Rows)-1].Version, nil
}

// Open ...
func (td *TDriver) Open(dsn string) error {
	return nil
}

// SetVersion ...
func (td *TDriver) SetVersion(version uint64) error {
	td.Rows = append(td.Rows, &Row{
		Version:   version,
		CreatedAt: time.Now(),
	})
	return nil
}

// DeleteVersion ...
func (td *TDriver) DeleteVersion(version uint64) error {
	var rs []*Row
	for _, r := range td.Rows {
		if r.Version != version {
			rs = append(rs, r)
		}
	}
	td.Rows = rs
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

// Versions returns only version columns
func (td *TDriver) Versions() []uint64 {
	var vs []uint64

	for _, r := range td.Rows {
		vs = append(vs, r.Version)
	}

	return vs
}
