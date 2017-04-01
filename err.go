package spirali

import "errors"

// Various errors the spirali might return.
var (
	ErrUnknownDriver         = errors.New("unknown driver")
	ErrEnvNotFound           = errors.New("env not found in config")
	ErrMigrationFileNotFound = errors.New("migration file not found")
	ErrMigrationsNotExist    = errors.New("migrations not exist")
	ErrSchemaVersionIsZero   = errors.New("schema version is 0")
)
