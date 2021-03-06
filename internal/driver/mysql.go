package driver

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" //use mysql
	"github.com/k0kubun/pp"
)

// Mysql is one of driver impls.
type Mysql struct {
	conn   *sql.DB
	tx     *sql.Tx
	locker sync.Mutex
}

// Close mysql connection.
func (m *Mysql) Close() error {
	return m.conn.Close()
}

// CreateVersionTableIfNotExists creates schema version table if table does not exist.
func (m *Mysql) CreateVersionTableIfNotExists() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		  id bigint UNSIGNED NOT NULL AUTO_INCREMENT,
		  version BIGINT UNSIGNED NOT NULL UNIQUE,
		  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`, schemaManagementTableName)

	if err := m.Exec(query); err != nil {
		return err
	}
	return nil
}

// DeleteVersion deletes version column from database.
func (m *Mysql) DeleteVersion(version uint64) error {
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE version = %d
	`, schemaManagementTableName, version)

	if err := m.Exec(query); err != nil {
		return err
	}
	return nil
}

// Exec executes query.
func (m *Mysql) Exec(query string) error {
	if m.tx != nil {
		if _, err := m.tx.Exec(query); err != nil {
			return err
		}
	} else {
		if _, err := m.conn.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

// GetAppliedTimeList returns list of migration applied times.
func (m *Mysql) GetAppliedTimeList() (map[uint64]time.Time, error) {
	query := fmt.Sprintf("SELECT version, created_at FROM %s order by version asc", schemaManagementTableName)

	data := map[uint64]time.Time{}

	rows, err := m.conn.Query(query)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var version uint64
		var createdAt time.Time

		if err := rows.Scan(&version, &createdAt); err != nil {
			return data, nil
		}

		data[version] = createdAt
	}

	return data, nil
}

// GetCurrentVersion returns current migration version of database.
func (m *Mysql) GetCurrentVersion() (uint64, error) {
	query := fmt.Sprintf("select version from %s order by version desc limit 1", schemaManagementTableName)

	var version uint64
	err := m.conn.QueryRow(query).Scan(&version)

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return version, nil
}

// Open mysql connection.
func (m *Mysql) Open(dsn string) error {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if err := conn.Ping(); err != nil {
		return err
	}
	m.conn = conn

	return nil
}

// SetVersion appends new version.
func (m *Mysql) SetVersion(version uint64) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (version) VALUES (%d)
  `, schemaManagementTableName, version)

	if err := m.Exec(query); err != nil {
		return err
	}
	return nil
}

// Transaction executes a function atomically.
// Roll back the process if function returns error.
func (m *Mysql) Transaction(fn func() error) error {
	m.locker.Lock()
	defer m.locker.Unlock()

	tx, err := m.conn.Begin()
	if err != nil {
		return err
	}
	m.tx = tx

	if err := fn(); err != nil {
		pp.Println(err)
		if e := m.tx.Rollback(); e != nil {
			return e
		}
		return err
	}
	if err := m.tx.Commit(); err != nil {
		if e := m.tx.Rollback(); e != nil {
			return e
		}
		return err
	}

	return nil
}
